package apierror

type Code string

const (
	CodeInvalidRequest   Code = "invalid_request"
	CodeAuthentication   Code = "unauthorized"
	CodePermission       Code = "permission_denied"
	CodeNotFound         Code = "not_found"
	CodeMethodNotAllowed Code = "method_not_allowed"
	CodeConflict         Code = "conflict"
	CodeRateLimit        Code = "rate_limited"
	CodeInternal         Code = "internal"
)

type Error struct {
	Status  int
	Code    Code
	Message string
	Hint    string
}

func New(status int, code Code, message string, hint string) error {
	return &Error{
		Status:  status,
		Code:    code,
		Message: message,
		Hint:    hint,
	}
}

func (e *Error) Error() string {
	return e.Message
}
