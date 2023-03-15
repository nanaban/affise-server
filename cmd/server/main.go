package main

import (
	"context"
	"log"

	"affise-server/internal/server"
)

func main() {
	s := server.NewServer(":8080") //todo

	if err := s.Run(context.Background()); err != nil {
		log.Fatalf("server run err: %v", err)
	}
}
