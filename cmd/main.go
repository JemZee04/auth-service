package main

import (
	auth_service "auth-service"
	"auth-service/pkg/handler"
	"log"
	"net/http"
)

func main() {
	handlers := new(handler.Handler)
	srv := new(auth_service.Server)
	if err := srv.Run("8080", http.HandlerFunc(handlers.Serve)); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
