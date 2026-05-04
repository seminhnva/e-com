package main

import (
	"log"
	"time"

	"github.com/seminhnva/e-com/internal/db"
	"github.com/seminhnva/e-com/internal/env"
	"github.com/seminhnva/e-com/internal/store"
)

const version = "1.0.0"

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:          env.GetString("DB_ADDR", "postgres://admin:secret@localhost:5432/e-com?sslmode=disable"),
			maxOpenConns:  env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIddleConns: env.GetInt("DB_MAX_IDDLE_CONNS", 25),
			maxIdleTime:   env.GetDuration("DB_MAX_IDLE_TIME", time.Minute*5),
		},
		env:     env.GetString("APP_ENV", "development"),
		version: version,
	}
	db, err := db.NewDB(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIddleConns, cfg.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Printf("db connection pool established")
	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
