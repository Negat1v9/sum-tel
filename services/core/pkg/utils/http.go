package utils

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Negat1v9/sum-tel/shared/logger"
)

const (
	ErrBadRequest        = "Bad request"
	ErrNoSuchUser        = "User not found"
	ErrNotFound          = "Not Found"
	ErrUnauthorized      = "Unauthorized"
	ErrForbidden         = "Forbidden"
	ErrBadQueryParams    = "Invalid query params"
	ErrRequestTimeout    = "Request timeout"
	ErrInternal          = "Internal server error"
	ErrNotEnoughStars    = "Not enough stars"
	ErrNotValidDuraction = "Not valid duration"
	ErrNoServerFound     = "No server found"
	ErrStarsAmount       = "Invalid stars amount minimum 150"
	ErrToManyRequests    = "Too many requests"
	ErrConflict          = "Conflict"
)

type Error struct {
	StatusCode int    `json:"status_code"` // http status code
	Message    string `json:"message"`     // error message
	Causes     any    `json:"-"`           // error causes for internal use
}

// Errors - implementation of the error interface
func (e *Error) Error() string {
	return fmt.Sprintf("status: %d, message: %s, causes: %v", e.StatusCode, e.Message, e.Causes)
}

// New creates a new HTTP error
func NewError(statusCode int, message string, causes any) *Error {
	return &Error{
		StatusCode: statusCode,
		Message:    message,
		Causes:     causes,
	}
}

// ParseError - parses an error into an HTTP error
func parseError(err error) *Error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return NewError(http.StatusNotFound, ErrNotFound, nil)
	case errors.Is(err, context.DeadlineExceeded):
		return NewError(http.StatusRequestTimeout, ErrRequestTimeout, nil)
	case strings.Contains(err.Error(), "Unmarshal"):
		return NewError(http.StatusBadRequest, ErrBadRequest, nil)
	default:
		if restErr, ok := err.(*Error); ok {
			return restErr
		}
		return NewError(http.StatusInternalServerError, ErrInternal, nil)
	}
}

func LogResponseErr(r *http.Request, log *logger.Logger, err error) {
	if err != nil {
		log.Errorf("Path: %s, Error: %s", r.RequestURI, err.Error())
	}
}

func WriteJsonResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func WriteErrResponse(w http.ResponseWriter, err error) {
	httpErr := parseError(err)

	WriteJsonResponse(w, httpErr.StatusCode, httpErr)
}
