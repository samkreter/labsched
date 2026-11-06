package main

import (
	"log"
	"net/http"

	"github.com/samkreter/labsched/internal/server"
)

func main() {
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", server.New()))
}
