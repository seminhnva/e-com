package main

import (
	"log"

	"github.com/seminhnva/e-com/internal/env"
)

func main() {
	cfg := Config{
		addr: env.GetString("ADDR", ":8080"),
	}
	app := &Application{
		config: cfg,
	}
	mux := app.mount()
	log.Fatal(app.run(mux))
}
