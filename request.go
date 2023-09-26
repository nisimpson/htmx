package htmx

import (
	"net/http"
	"net/url"
)

type Request struct {
	*http.Request
}

func NewRequest(r *http.Request) *Request {
	return &Request{Request: r}
}

// IsHTMXRequest returns true if the "HX-Request" header key has a value of "true".
// This indicates that the client is using the htmx library to render html.
func (r Request) IsHTMXRequest() bool {
	return r.Header.Get(HeaderHXRequest) == "true"
}

// HTMXTriggerName returns the id of the element that triggered the
// server request. This value is stored within the "HX-Trigger" header.
func (r Request) HTMXTriggerID() string {
	return r.Header.Get(HeaderHXTrigger)
}

// HTMXTriggerName returns the name of the element that triggered the
// server request. This value is stored within the "HX-Trigger-Name" header.
func (r Request) HTMXTriggerName() string {
	return r.Header.Get(HeaderHXTriggerName)
}

// HTMXTargetID returns the id of the target element receiving the html
// fragment response. This value is stored within the "HX-Target" header.
func (r Request) HTMXTargetID() string {
	return r.Header.Get(HeaderHXTarget)
}

// HTMXPrompt returns the value entered by the client user when
// prompted via the "hx-prompt" attribute. This value is stored within
// the "HX-Prompt" header.
func (r Request) HTMXPrompt() string {
	return r.Header.Get(HeaderHXPrompt)
}

// HTMXCurrentURL returns the current URL of the client browser if specified.
func (r Request) HTMXCurrentURL() (*url.URL, bool) {
	if value, err := url.Parse(r.Header.Get(HeaderHXCurrentURL)); err != nil {
		return nil, false
	} else {
		return value, true
	}
}

// IsHTMXHistoryRestoreRequest returns true if the "HX-History-Restore-Request"
// header is set to "true". This indicates that this request is for history restoration
// after a miss in the local history cache.
func (r Request) IsHTMXHistoryRestoreRequest() bool {
	return r.Header.Get(HeaderHXHistoryRestoreRequest) == "true"
}
