package main

import (
	"log"
	"net/http"
	"os"

	"probable-invention/pkg/api"
	"probable-invention/pkg/db"
)

const defaultPort = "7540"

func main() {
	dbFile := "scheduler.db"
	if v := os.Getenv("TODO_DBFILE"); v != "" {
		dbFile = v
	}
	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}
	api.Init()

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	webDir := "web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	addr := ":" + port
	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
