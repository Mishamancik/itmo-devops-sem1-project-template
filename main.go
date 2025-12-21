package main

import (
	"log"

	"project_sem/internal/db"
	"project_sem/internal/server"
)

func main() {
	database, err := db.New()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	srv, err := server.New(database)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting server on :8080")
	if err := srv.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
