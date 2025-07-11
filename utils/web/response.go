package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

const (
	ErrValidation     = 1001
	ErrAuthentication = 1002
	ErrPermission     = 1003
	ErrNotFound       = 1004
	ErrConflict       = 1005
	ErrProcessing     = 1500
)

// Response adalah struktur standar untuk semua response API
type Response struct {
	Message    string              `json:"message,omitempty" extensions:"x-omitempty,x-nullable"` // Pesan umum, bisa kosong
	Status     string              `json:"status"`
	Data       any                 `json:"data,omitempty" extensions:"x-omitempty,x-nullable"`
	ErrorCode  int                 `json:"code,omitempty" extensions:"x-omitempty,x-nullable"`
	Pagination *PaginationResponse `json:"pagination,omitempty" extensions:"x-omitempty,x-nullable"`
}

func (r *Response) SetData(data any) {
	r.Data = data
}

// writeJSON sets headers and writes v as JSON with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func OK(w http.ResponseWriter, status int, data any, pagination ...PaginationResponse) error {
	resp := Response{
		Status: "success",
		Data:   data,
	}
	if len(pagination) > 0 {
		resp.Pagination = &pagination[0]
	}
	return writeJSON(w, status, resp)
}

func OKNoContent(w http.ResponseWriter, status int) error {
	resp := Response{
		Status: "success",
	}

	return writeJSON(w, status, resp)
}

// Err inspects err for validation or HTTP errors, then writes a JSON error envelope
func Err(w http.ResponseWriter, err error) error {
	slog.Info("", "error", err)

	// 1. Validation errors
	var vErrs validator.ValidationErrors
	if errors.As(err, &vErrs) {
		var msgs string
		for _, fe := range vErrs {
			msgs = fmt.Sprintf("Field '%s' failed on the '%s' rule", fe.Field(), fe.Tag())
		}
		return writeJSON(w, http.StatusBadRequest, Response{
			Status:    "failed",
			Message:   msgs,
			ErrorCode: ErrValidation,
		})
	}

	// 2. a custom *HTTPError type
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return writeJSON(w, httpErr.Code, Response{
			Status:    "error",
			Message:   httpErr.Message,
			ErrorCode: httpErr.ErrorCode,
		})
	}

	// 3. Fallback: 500 Internal Server Error
	slog.Error("Unhandled internal error", "error", err)
	return writeJSON(w, http.StatusInternalServerError, Response{
		Status: "error",
		// Errors: []string{"An internal server error occurred"},
		Message:   err.Error(),
		ErrorCode: ErrProcessing,
	})
}

// Example of a custom HTTPError you might use in your handlers:
type HTTPError struct {
	Code      int
	Message   string
	ErrorCode int
}

func (e *HTTPError) Error() string {
	return e.Message
}

func NewHTTPError(code int, message string, errorCode int) *HTTPError {
	return &HTTPError{
		Code:      code,
		Message:   message,
		ErrorCode: errorCode, // Default error code, can be set later
	}
}
