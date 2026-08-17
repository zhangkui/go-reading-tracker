package main

import (
	"log"
	"net/http"
	"os"

	"github.com/zhangkui/go-reading-tracker/internal/httpapi"
	"github.com/zhangkui/go-reading-tracker/internal/reading"
)

func main() {
	address := os.Getenv("ADDR")
	if address == "" {
		address = ":8080"
	}
	server := httpapi.NewServer(reading.NewService())
	log.Printf("reading tracker listening on %s", address)
	log.Fatal(http.ListenAndServe(address, server.Handler()))
}
