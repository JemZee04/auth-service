package handler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	logrus.Error(message)

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse{message})
}
