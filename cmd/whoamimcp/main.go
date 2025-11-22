package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/traefik/whoamimcp/internal"
)

var (
	addr string
	name string
)

func init() {
	name = os.Getenv("WHOAMI_NAME")
	if name == "" {
		name = "whoamimcp"
	}
	host := os.Getenv("WHOAMI_HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	port := os.Getenv("WHOAMI_PORT_NUMBER")
	if port == "" {
		port = "80"
	}
	flag.StringVar(&name, "name", name, "give me a name")
	flag.StringVar(&addr, "addr", fmt.Sprintf("%s:%s", host, port), "http service address")
}

func main() {
	flag.Parse()
	ctx := context.Background()
	log.Printf("Listening at: http://%s", addr)
	log.Fatal(internal.StartServer(ctx, name, addr))
}
