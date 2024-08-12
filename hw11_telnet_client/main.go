package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

const defaultTimeout = 10

func main() {
	// parse command window arguments
	timeout := flag.Duration("timeout", defaultTimeout*time.Second, "timeout connection default 10s")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("Port or host not specified")
	}

	host := args[0]
	port := args[1]
	log.Printf("Host: %v, Port: %v, Timeout: %v", host, port, timeout)

	// new client
	clientTel := NewTelnetClient(net.JoinHostPort(host, port), *timeout, os.Stdin, os.Stdout)

	// run
	if err := runTelnet(clientTel); err != nil {
		log.Printf("failed with error: %s", err)
		os.Exit(1)
	}
}

func runTelnet(clientTel TelnetClient) error {
	// connect
	if err := clientTel.Connect(); err != nil {
		return fmt.Errorf("connection error %w", err)
	}
	defer func(clientTel TelnetClient) {
		err := clientTel.Close()
		if err != nil {
			log.Fatalf("Error while close connection: %s", err)
		}
	}(clientTel)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	go func() {
		defer stop()
		if err := clientTel.Send(); err != nil {
			log.Fatalf("Error reading from channel. Error: %s", err)
		}
	}()

	go func() {
		defer stop()
		if err := clientTel.Receive(); err != nil {
			log.Fatalf("Error reading from channel. Error: %s", err)
		}
	}()

	<-ctx.Done()
	return nil
}
