package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/seminhnva/e-com/internal/store"
)

type application struct {
	config config
	store  store.Storage
	db     dbConfig
}

type config struct {
	addr    string
	db      dbConfig
	env     string
	version string
}
type dbConfig struct {
	addr          string
	maxOpenConns  int
	maxIddleConns int
	maxIdleTime   time.Duration
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(1 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		//v1/posts
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Patch("/", app.updatePostHandler)
			})
		},
		)
	})
	return r
}
func (app *application) run(mux *chi.Mux) error {
	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	log.Printf("server has started at %s", app.config.addr)
	return srv.ListenAndServe()
}
