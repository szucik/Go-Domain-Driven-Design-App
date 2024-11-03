package apperrors

import (
	"net/http"

	"encoding/json"

	"github.com/google/uuid"
)

type ErrorResponse struct {
	TraceId string `json:"traceId"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

func (e ErrorResponse) JSONError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(e.Code)
	json.NewEncoder(w).Encode(e)
}

func Error(message, errorType string, code int) ErrorResponse {
	return ErrorResponse{
		TraceId: uuid.New().String(),
		Code:    code,
		Message: message,
		Type:    errorType,
	}
}
