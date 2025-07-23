package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Инициализация базы данных
	if err := InitDB(); err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	defer CloseDB()

	// Настройка маршрутов
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/products", GetProducts).Methods("GET")
	api.HandleFunc("/products", CreateProduct).Methods("POST")
	api.HandleFunc("/products/{id:[0-9]+}", GetProduct).Methods("GET")
	api.HandleFunc("/products/{id:[0-9]+}", UpdateProduct).Methods("PUT")
	api.HandleFunc("/products/{id:[0-9]+}", DeleteProduct).Methods("DELETE")

	// Middleware
	r.Use(loggingMiddleware)

	// Запуск сервера
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
