package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/seminhnva/e-com/docs"
	"github.com/seminhnva/e-com/internal/store"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config config
	store  store.Storage
	db     dbConfig
	logger *zap.SugaredLogger
}

type config struct {
	addr    string
	db      dbConfig
	env     string
	version string
	apiURL  string
	mail    mailConfig
}
type mailConfig struct {
	exp time.Duration
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
	r.Use(middleware.Timeout(3 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		//v1/posts
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Patch("/", app.updatePostHandler)
				r.Route("/comments", func(r chi.Router) {
					r.Post("/", app.createCommentHandler)
				})
			})
		},
		)

		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.userByIdContextMiddleware)
				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})
		})

		r.Group(func(r chi.Router) {
			r.Get("/feed", app.getFeedHandler)

		})
		r.Route("authentication", func(r chi.Router) {
			// r.Post("/user", app.loginHandler)
			r.Post("/user", app.registerUserHandler)
		})
	})
	return r
}
func (app *application) run(mux *chi.Mux) error {
	//Docs
	docs.SwaggerInfo.Version = app.config.version
	if u, err := url.Parse(app.config.apiURL); err == nil {
		docs.SwaggerInfo.Host = u.Host
	}
	docs.SwaggerInfo.BasePath = "/v1"
	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	app.logger.Infof("server has started at %s", app.config.addr)
	return srv.ListenAndServe()
}
