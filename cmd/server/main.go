package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

//go:generate protoc -I ./../.. --go_out ./../.. ./../../pkg/data/*.proto
func main() {
	_ = godotenv.Load()

	router := NewRouter()

	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = ":3000"
	}

	log.Printf("Listening on %s", listen)
	panic(router.Run(listen))
}
