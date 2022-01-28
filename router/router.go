package router

import (
    "acad-be/middleware"
    "github.com/gorilla/mux"
)

func Router() *mux.Router {
    router := mux.NewRouter()
    router.HandleFunc("/api/create-user", middleware.CreateUser).Methods("POST", "OPTIONS")
    router.HandleFunc("/api/user/{id}", middleware.GetUser).Methods("GET", "OPTIONS")

    return router
}