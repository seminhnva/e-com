package main

import (
	"log"
	"time"

	"github.com/seminhnva/e-com/internal/db"
	"github.com/seminhnva/e-com/internal/env"
	"github.com/seminhnva/e-com/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:secret@localhost:5432/e-com?sslmode=disable")

	conn, err := db.NewDB(addr, 3, 3, 15*time.Minute)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	store := store.NewStorage(conn)
	db.Seed(store, conn)
}
