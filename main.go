package main

import (
	"flag"
	"log"
	"os"
	"time"

	"go.fodro/nyx/config"
	"go.fodro/nyx/http"
)

func main() {
	flag.Parse()
	c, err := config.Load()
	if err != nil {
		log.Printf("Could not read configuration: %s\n", err)
		return
	}

	log.Printf("Starting Server at %s\n", c.ListenOn)
	if err := http.Start(c); err != nil {
		log.Printf("Could not start server or server crashed: %s\n", err)
		log.Printf("Waiting 10 seconds before dying...")
		time.Sleep(10 * time.Second)
		log.Printf("Exiting")
		os.Exit(1)
		return
	}
	os.Exit(0)
}
