package htmx

import "net/http"

// Handler instances respond to htmx requests, which are just HTML requests with
// some special headers that provide context for htmx-centric events. To help server developers
// adhere to client expectations, the respective writer and request arguments provide additional
// methods to the standard library structs of the same name.
type Handler interface {
	// ServeHTMX is invoked when the server is responding to an incoming http request with htmx
	// context.
	ServeHTMX(*ResponseWriter, *Request)
}

// The HandlerFunc type is an adapter to allow the use of ordinary functions as htmx handlers.
// If f is a function with the appropriate signature, HandlerFunc(f) is a Handler that calls f.
type HandlerFunc func(*ResponseWriter, *Request)

// ServeHTMX calls f(w, r).
func (f HandlerFunc) ServeHTMX(w *ResponseWriter, r *Request) {
	f(w, r)
}

// HTMX wraps the htmx handler into a standard library http handler function,
// which can be used by a Go http muxer.
func HTMX(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTMX(NewResponseWriter(w), NewRequest(r))
	}
}

// HTMXFunc wraps the htmx handler function into a standard library http handler function,
// which can be used by a Go http muxer.
func HTMXFunc(f HandlerFunc) http.HandlerFunc {
	return HTMX(f)
}
