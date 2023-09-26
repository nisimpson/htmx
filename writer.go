package htmx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// ResponseWriter is responsible for write htmx compliant http responses
// back to the client.
type ResponseWriter struct {
	http.ResponseWriter
}

// NewResponseWriter creates a new htmx response writer instance,
// using the underlying http response writer to generate a http response.
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w}
}

// SetPushHeader sets the "HX-Push" header which triggers the web client to
// push the target URL into the browser's address bar.
func (r ResponseWriter) SetPushHeader(url url.URL) {
	r.Header().Set(HeaderHXPushURL, url.String())
}

// SetRedirectHeader sets the "HX-Redirect" header which triggers the web client
// to redirect to a new URL.
func (r ResponseWriter) SetRedirectHeader(url url.URL) {
	r.Header().Set(HeaderHXRedirect, url.String())
}

// SetLocationHeader sets the "HX-Location" header which triggers the web client
// to redirect to a new URL that acts as a swap.
func (r ResponseWriter) SetLocationHeader(url url.URL) {
	r.Header().Set(HeaderHXLocation, url.String())
}

// SetRefreshHeader sets the "HX-Refresh" header which triggers a full refresh
// of the web client.
func (r ResponseWriter) SetRefreshHeader() {
	r.Header().Set(HeaderHXRefresh, "true")
}

// TriggerEvent defines a single event, or multiple events that should be
// triggered client side once a htmx response is received.
type TriggerEvent interface {
	triggerHeaderValue() string
}

// SetTriggerHeader writes the "HX-Trigger-After-Swap" header
// to the http response, triggering client side event(s) upon receipt.
func (r ResponseWriter) SetTriggerHeader(event TriggerEvent) {
	r.Header().Set(HeaderHXTrigger, event.triggerHeaderValue())
}

// SetTriggerAfterSettleHeader writes the "HX-Trigger-After-Swap" header
// to the http response, triggering client side event(s) after the settling step.
func (r ResponseWriter) SetTriggerAfterSettleHeader(event TriggerEvent) {
	r.Header().Set(HeaderHXTriggerAfterSettle, event.triggerHeaderValue())
}

// SetTriggerAfterSwapHeader writes the "HX-Trigger-After-Swap" header
// to the http response, triggering client side event(s) after the swap step.
func (r ResponseWriter) SetTriggerAfterSwapHeader(event TriggerEvent) {
	r.Header().Set(HeaderHXTriggerAfterSwap, event.triggerHeaderValue())
}

type Component interface {
	RenderHTMX(io.Writer) error
}

type ComponentFunc func(io.Writer) error

func (f ComponentFunc) RenderHTMX(w io.Writer) error {
	return f(w)
}

// WriteComponent invokes the Render() method on the provided component,
// writing the contents to the http response writer.
func WriteComponent(w http.ResponseWriter, component Component, status int) {
	// initialize new buffer; a temporary buffer is used to ensure
	// the template transformation is valid and safe to transport back
	// to the client.
	buf := bytes.Buffer{}

	err := component.RenderHTMX(&buf)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

// TriggerEvents creates an event trigger for either a single event or multiple
// events, without providing any additional details. To provide JSON context
// with events, use the TriggerEventsWithDetails() function.
func TriggerEvents(eventName string, eventNames ...string) triggerEvents {
	eventNames = append(eventNames, eventName)
	return triggerEvents{eventNames: eventNames}
}

// TriggerEventsWithContext creates an event trigger for either a single event or multiple
// events, providing event details via JSON. For example:
//
//	w.SetTriggerHeader(htmx.TriggerEventsWithContext(
//		"showMessage": "Here is a message"
//	))
//
// triggers a single event "showMessage" with value "Here is a message". To pass
// multiple pieces of data use a map or JSON encodable struct:
//
//	w.SetTriggerHeader(htmx.TriggerEventsWithContext(
//		"showMessage": map[string]any{
//			"level": "info",
//			"message": "Here is a message"
//		}
//	))
//
// To send multiple events, add another key value pair to the eventsData parameter:
//
//	w.SetTriggerHeader(htmx.TriggerEventsWithContext(
//		"event1": map[string]any{
//			"level": "info",
//			"message": "Here is a message"
//		},
//		"event2": map[string]any{
//			"level": "info",
//			"message": "Here is another message"
//		},
//	))
//
// If the provided context cannot be encoded into JSON, this function panics.
func TriggerEventsWithContext(context map[string]any) triggerEventsJSON {
	return triggerEventsJSON{eventContext: context}
}

type triggerEvents struct {
	eventNames []string
}

func (s triggerEvents) triggerHeaderValue() string {
	return strings.Join(s.eventNames, ",")
}

type triggerEventsJSON struct {
	eventContext map[string]any
}

func (m triggerEventsJSON) triggerHeaderValue() string {
	data, err := json.Marshal(m.eventContext)
	if err != nil {
		err = fmt.Errorf("failed to marshal trigger message event: %s", err)
		panic(err)
	}
	return string(data)
}
