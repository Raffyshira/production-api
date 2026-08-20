package main

import (
	"log"

	"github.com/raffyshira/project-rest-api/internal/db"
	"github.com/raffyshira/project-rest-api/internal/env"
	"github.com/raffyshira/project-rest-api/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/dbname?sslmode=disable")

	conn, err := db.New(addr, 3, 3, "15m")

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store, conn)
}
