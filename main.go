package main

import (
	"acad-be/router"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {
  r := router.Router()
  credentials := handlers.AllowCredentials()
  methods := handlers.AllowedMethods([]string{"POST", "GET"})
  origins := handlers.AllowedOrigins([]string{"*"})

  fmt.Println("Starting server on the port 8080...")

  log.Fatal(http.ListenAndServe(":8080", handlers.CORS(credentials, methods, origins)(r)))
}