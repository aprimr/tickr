package response

import (
	"encoding/json"
	"net/http"
)

// JSON sends a success response with given status code, message and data
func JSON(w http.ResponseWriter, statusCode int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := SuccessResponse{
		Success: true,
		Status:  statusCode,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

// Error sends a error response with given status code, message and errors
func Error(w http.ResponseWriter, statusCode int, message string, errors any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Success: false,
		Status:  statusCode,
		Message: message,
		Errors:  errors,
	}

	json.NewEncoder(w).Encode(response)
}
