package main

import (
	"fmt"
	"log"
	"net/http"
	"rickandmorty-backend/routes"
)

func main() {
	http.HandleFunc("/search", routes.SearchHandler)
	http.HandleFunc("/top-pairs", routes.TopPairsHandler)

	port := ":8080"
	fmt.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
