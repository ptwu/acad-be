package main

import (
	"acad-be/router"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func determineListenAddress() (string, error) {
  port := os.Getenv("PORT")
  if port == "" {
    return "", fmt.Errorf("$PORT not set")
  }
  return ":" + port, nil
}

func main() {
  r := router.Router()
  credentials := handlers.AllowCredentials()
  methods := handlers.AllowedMethods([]string{"POST", "GET"})
  origins := handlers.AllowedOrigins([]string{"*"})

  fmt.Println("Starting server on the port 8080...")
  addr, err := determineListenAddress()
  if err != nil {
    log.Fatal(err)
  }
  log.Fatal(http.ListenAndServe(addr, handlers.CORS(credentials, methods, origins)(r)))
}