package main

import (
	"log"

	"github.com/balajiss36/go-webrtc/internal/server"
)

func main() {
	err := server.Run()
	if err != nil {
		log.Fatalf("Error in running server: %v", err)
	}
}
