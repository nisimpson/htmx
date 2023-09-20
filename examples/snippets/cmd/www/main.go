package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nisimpson/htmx"
	"github.com/nisimpson/htmx/examples/snippets"
	"github.com/nisimpson/htmx/examples/snippets/pkg/models"
	"github.com/nisimpson/htmx/examples/snippets/pkg/storage"
)

func main() {
	app := SnippetBox{
		SnippetModel: models.SnippetModel{
			Store: storage.NewMemoryStorage(),
		},
		Templates: initTemplateCache(),
	}
	err := http.ListenAndServe(":3333", &app)
	log.Fatalln(err)
}

type SnippetBox struct {
	models.SnippetModel
	Templates map[string]*template.Template
}

func (s *SnippetBox) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	s.logJSON(r.Header)

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("./assets/"))

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/assets/". For matching paths, we strip the
	// "/assets" prefix before the request reaches the file server.
	mux.Handle("/assets/", http.StripPrefix("/assets", s.noIndex(fileServer)))

	// Register the other application routes as normal.
	mux.HandleFunc("/", htmx.HTMXFunc(s.home))
	mux.HandleFunc("/snippet", htmx.HTMXFunc(s.viewSnippet))
	mux.HandleFunc("/snippets/create", htmx.HTMXFunc(s.createSnippet))
	mux.HandleFunc("/snippets", htmx.HTMXFunc(s.snippets))

	// protect from cross site scripting and serve.
	s.secureHeaders(mux).ServeHTTP(w, r)
}

type SnippetView interface {
	TemplateName() string
}

type TemplateFilesProvider interface {
	TemplateFiles() []string
}

type TemplateDataProvider interface {
	TemplateData() any
}

var functions = template.FuncMap{
	"currentYear": currentYear,
	"humanDate":   humanDate,
}

func (s SnippetBox) render(w *htmx.ResponseWriter, r *htmx.Request, view SnippetView) {
	var fn htmx.ComponentFunc = func(ctx context.Context, w io.Writer) error {
		// retrieve the approperiate template set from the cache based on the template name
		name := view.TemplateName()
		ts, ok := s.Templates[name]

		// if no entry exists in the cache; query the view for the relevant template
		// files to be parsed.
		if !ok {
			provider, ok := view.(TemplateFilesProvider)
			if !ok {
				err := fmt.Errorf("the template %s does not exist", name)
				log.Println(err)
				return err
			}

			parsed, err := template.New(name).
				Funcs(functions).
				ParseFiles(provider.TemplateFiles()...)

			if err != nil {
				log.Print(err.Error())
				return err
			}
			ts = parsed
		}

		var data any = view
		if observer, ok := view.(TemplateDataProvider); ok {
			// allow view to make any changes before template is executed
			data = observer.TemplateData()
		}

		// Use the ExecuteTemplate() method to write the content of the
		// template as the response body.
		err := ts.ExecuteTemplate(w, view.TemplateName(), data)
		if err != nil {
			log.Print(err.Error())
			return err
		}

		return nil
	}

	w.WriteHeader(http.StatusOK)
	w.WriteComponent(r.Context(), fn)
}

func (SnippetBox) serverError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (SnippetBox) noIndex(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (SnippetBox) secureHeaders(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1;mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	}
}

func (SnippetBox) logJSON(data any) {
	out, _ := json.MarshalIndent(data, "", "  ")
	log.Println(string(out))
}

func initTemplateCache() map[string]*template.Template {
	cache := make(map[string]*template.Template)
	rootDir := snippets.RootDir()

	pages, err := filepath.Glob(filepath.Join(rootDir, "html/pages/*.tmpl"))
	if err != nil {
		panic(err)
	}

	for _, page := range pages {
		// extract the file name (ex. home.tmpl) from file path
		name := filepath.Base(page)

		// parse page template file into a template set
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			panic(err)
		}

		// add any component templates to the set
		ts, err = ts.ParseGlob(filepath.Join(rootDir, "html/components/*.tmpl"))
		if err != nil {
			panic(err)
		}

		// add template to the cache
		log.Println("adding template to cache:", name)
		cache[name] = ts
	}
	return cache
}

func currentYear() string {
	return strconv.Itoa(time.Now().Year())
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}
