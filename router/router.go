package router

import (
    "acad-be/middleware"
    "github.com/gorilla/mux"
)

func Router() *mux.Router {
    router := mux.NewRouter()
    router.HandleFunc("/api/create-user", middleware.CreateUser).Methods("POST", "OPTIONS")
    router.HandleFunc("/api/user/{id}", middleware.GetUser).Methods("GET", "OPTIONS")
		router.HandleFunc("/api/switch-basis/{id}", middleware.SwitchBasis).Methods("POST", "OPTIONS")
		router.HandleFunc("/api/review-card/{id}", middleware.ReviewCard).Methods("POST", "OPTIONS")
    return router
}