package httputils

import (
	"encoding/json"
	"net/http"
)

func DispatchNewHttpError(w http.ResponseWriter, message string, statusCode int) {
	responseContent, _ := json.Marshal(map[string]string{"message": message})

	w.WriteHeader(statusCode)
	w.Write(responseContent)
}

func DispatchNewResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	responseContent, _ := json.Marshal(data)

	w.WriteHeader(statusCode)
	w.Write(responseContent)
}

func WrapAPIResponse(data interface{}, message string) map[string]interface{} {
	return map[string]interface{}{
		"data":    data,
		"message": message,
	}
}
