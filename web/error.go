package web

import (
	"encoding/json"
	"net/http"

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

func BadRequestError(message, errorType string) ErrorResponse {
	return ErrorResponse{
		TraceId: uuid.New().String(),
		Code:    400,
		Message: message,
		Type:    errorType,
	}
}
