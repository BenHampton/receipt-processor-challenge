package main

import (
	"fetch-interview/internal/routes"
	"fmt"
	"net/http"
)

func main() {
	routers := routes.NewRouter()
	port := 8080
	address := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting server on port %s\n", address)
	err := http.ListenAndServe(address, routers)
	if err != nil {
		panic(err)
	}
}
