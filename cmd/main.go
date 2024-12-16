package main

import (
	"github.com/balajiss36/go-webrtc/internal/server"
	"github.com/rs/zerolog/log"
)

func main() {
	err := server.Run()
	if err != nil {
		log.Fatalf("Error in running server: %v", err)
	}
}
