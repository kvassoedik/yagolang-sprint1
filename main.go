package main

import (
	"log"
	"net/http"

	"final/handler"
)

func main() {
	http.HandleFunc("/api/v1/calculate", handler.CalculateHandler)

	log.Println("Server is running on http://127.0.0.1:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
