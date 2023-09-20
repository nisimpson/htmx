package htmx

const (
	// Request header. Always set to "true".
	//   - https://htmx.org/docs/#requests
	HeaderHXRequest = "HX-Request"
	// If defined within a request, contains the id of the element that=
	// triggered the request. If this header is defined within a response,
	// it will trigger client side events.
	//   - https://htmx.org/docs/#requests
	//   - https://htmx.org/headers/hx-trigger/
	HeaderHXTrigger = "HX-Trigger"
	// Request header containing the name of the element that triggered the request.
	//   - https://htmx.org/docs/#requests
	HeaderHXTriggerName = "HX-TriggerName"
	// Request header containing the id of the target element.
	//   - https://htmx.org/docs/#requests
	HeaderHXTarget = "HX-Target"
	// Request header containing the value entered by the user when prompted by
	// the "hx-prompt" attribute on the request element.
	//   - https://htmx.org/docs/#requests
	//   - https://htmx.org/attributes/hx-prompt/
	HeaderHXPrompt = "HX-Prompt"
	// Repsonse header that pushes a new URL into the browser's address bar.
	//   - https://htmx.org/docs/#requests
	HeaderHXPush = "HX-Push"
	// Response header that triggers a client-side redirect to a new location.
	//   - https://htmx.org/docs/#requests
	HeaderHXRedirect = "HX-Redirect"
	// Response header that triggers a client-side redirect to a new location
	// that acts as a swap.
	//   - https://htmx.org/docs/#requests
	HeaderHXLocation = "HX-Location"
	// Response header that when set to "true" will perform a full client-side
	// refresh.
	//   - https://htmx.org/docs/#requests
	HeaderHXRefresh = "HX-Refresh"
	// Response header that will trigger the named client-side event(s) after the
	// swap step.
	//   - https://htmx.org/headers/hx-trigger/
	HeaderHXTriggerAfterSwap = "HX-Trigger-After-Swap"
	// Response header that will trigger the named client-side event(s) after the
	// settle step.
	//   - https://htmx.org/headers/hx-trigger/
	HeaderHXTriggerAfterSettle = "HX-Trigger-After-Settle"
)
