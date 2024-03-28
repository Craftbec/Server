package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Craftbec/Server/internal/httpserver"
	"github.com/Craftbec/Server/internal/storage"
)

func main() {
	storeData, err := storage.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-stop
		cancel()
		close(done)
	}()
	go func() {
		err := httpserver.HTTPServer(ctx, storeData)
		if err != nil {
			log.Fatalf("HTTP server error: %v\n", err)
		}
	}()
	<-done
	storeData.GracefulStopDB()
	log.Println("Shutting down gracefully")
}
