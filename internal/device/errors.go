package device

import "net/http"

// DomainError represents a business-rule violation you can map to HTTP.
type DomainError struct {
	Code    string // e.g. "immutable_field", "forbidden_change"
	Field   string // e.g. "creation_time", "name"
	Message string // human-friendly
	HTTP    int    // suggested HTTP status (transport hint)
}

// Ensure *DomainError implements error.
var _ error = (*DomainError)(nil)

func (e *DomainError) Error() string { return e.Message }

func ErrImmutable(field string) *DomainError {
	return &DomainError{
		Code:    "immutable_field",
		Field:   field,
		Message: field + " is immutable",
		HTTP:    http.StatusBadRequest, // 400
	}
}

func ErrForbiddenChange(field, reason string, httpStatus int) *DomainError {
	return &DomainError{
		Code:    "forbidden_change",
		Field:   field,
		Message: "cannot change " + field + ": " + reason,
		HTTP:    httpStatus,
	}
}

func ErrRequired(field string) *DomainError {
	return &DomainError{
		Code: "required", Field: field,
		Message: field + " is required",
		HTTP:    http.StatusBadRequest,
	}
}

func ErrInvalid(field, reason string, httpStatus int) *DomainError {
	return &DomainError{
		Code: "invalid_" + field, Field: field,
		Message: reason,
		HTTP:    httpStatus,
	}
}
