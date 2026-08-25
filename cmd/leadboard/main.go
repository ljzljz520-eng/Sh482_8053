package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"enterpriselead/internal/service"
	"enterpriselead/internal/storage"
	"enterpriselead/internal/transport"
)

func main() {
	databasePath := flag.String("db", "leadboard.db", "path to bbolt database")
	address := flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()
	store, err := storage.Open(*databasePath)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	clock := service.FixedTime{Value: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
	app := service.New(store, clock)
	server := &http.Server{Addr: *address, Handler: transport.New(app).Handler()}
	log.Printf("enterprise procurement leadboard listening on %s", *address)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
