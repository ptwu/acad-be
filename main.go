package main

import (
	"acad-be/router"
	"fmt"
	"log"
	"net/http"
	"os"
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
  // credentials := handlers.AllowCredentials()
  // methods := handlers.AllowedMethods([]string{"POST", "GET"})
  // origins := handlers.AllowedOrigins([]string{"*"})

  addr, err := determineListenAddress()
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Starting server on the port " + addr)
  log.Fatal(http.ListenAndServe(addr, r))
}