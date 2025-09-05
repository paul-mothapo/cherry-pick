package main

import (
	"flag"
	"log"

	"github.com/cherry-pick/pkg/api"
)

func main() {
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

	log.Printf("Starting Database Intelligence Server on port %s", *port)
	log.Printf("UI available at: http://localhost:%s", *port)
	log.Printf("API available at: http://localhost:%s/api", *port)

	server := api.NewServer(*port)
	if err := server.Run(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
