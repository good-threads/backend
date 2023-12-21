package handlers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func respondMessage(w http.ResponseWriter, httpStatus int, message string) {
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(Response{
		Message: message,
	})
}
