package main

import (
	"context"
	"flag"
	"log"

	"affise-server/internal/config"
	"affise-server/internal/server"
)

var (
	flagServerAddr = flag.String("addr", ":8080", "server address")
)

func main() {
	flag.Parse()

	conf := config.NewDefault()
	conf.Server.Addr = *flagServerAddr

	s := server.NewServer(conf)

	if err := s.Run(context.Background()); err != nil {
		log.Fatalf("server run err: %v", err)
	}
}
