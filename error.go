package bitgo

// The error types are based on HTTP status codes.
const (
	// ErrorTypeRequiresApproval indicates that request is accepted but requires approval.
	ErrorTypeRequiresApproval = "requires_approval"
	// ErrorTypeInvalidRequest errors arise when a request has invalid parameters.
	ErrorTypeInvalidRequest = "invalid_request_error"
	// ErrorTypeAuthentication is returned when a request is not authenticated.
	ErrorTypeAuthentication = "authentication_error"
	// ErrorTypeNotFound is returned when API resource is not found.
	ErrorTypeNotFound = "not_found"
	// ErrorTypeRateLimit indicates too many requests hit the API too quickly.
	ErrorTypeRateLimit = "rate_limit_error"
	// ErrorTypeAPI covers temporary problems with BitGo's API (50x status codes).
	ErrorTypeAPI = "api_error"
)

// Error is the response returned when a call is unsuccessful.
type Error struct {
	// Type is an API error type based on HTTP status code.
	Type           string
	HTTPStatusCode int
	// Body is the raw response returned by the server.
	Body      string
	Message   string `json:"error"`
	RequestID string `json:"requestId"`
}

func (e Error) Error() string {
	return e.Message
}

// IsApprovalRequired returns true if err indicates that request is accepted but requires approval.
func (e Error) IsApprovalRequired() bool {
	return e.Type == ErrorTypeRequiresApproval
}

// IsInvalidRequest returns true if err caused by invalid request parameters.
func (e Error) IsInvalidRequest() bool {
	return e.Type == ErrorTypeInvalidRequest
}

// IsUnauthorized returns true if err caused by authentication problem.
func (e Error) IsUnauthorized() bool {
	return e.Type == ErrorTypeAuthentication
}

// IsNotFound returns true if err indicates that API resource is not found.
func (e Error) IsNotFound() bool {
	return e.Type == ErrorTypeNotFound
}

// IsRateLimited returns true if err caused by API requests throttling.
func (e Error) IsRateLimited() bool {
	return e.Type == ErrorTypeRateLimit
}

// IsTemporary returns true if err is temporary API error (50x status codes).
func (e Error) IsTemporary() bool {
	return e.Type == ErrorTypeAPI
}
