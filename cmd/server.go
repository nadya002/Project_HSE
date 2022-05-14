package main

import (
	"flag"
	"hse/pkg/server"
	"log"
)

func main() {
	host := flag.String("host", "127.0.0.1", "listen host")
	port := flag.Int("port", 7777, "listen port")
	numberOfClients := flag.Int("clients", 2, "number of input streams to merge")
	flag.Parse()
	if err := server.NewRunner(server.RunnerConfig{
		Host:            *host,
		Port:            *port,
		NumberOfClients: *numberOfClients,
	}).Run(); err != nil {
		log.Fatalln(err)
	}
}
