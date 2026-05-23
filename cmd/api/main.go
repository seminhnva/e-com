package main

import (
	"time"

	"github.com/seminhnva/e-com/internal/db"
	"github.com/seminhnva/e-com/internal/env"
	"github.com/seminhnva/e-com/internal/mailer"
	"github.com/seminhnva/e-com/internal/store"
	"go.uber.org/zap"
)

//	@title			Go social API
//	@description	An API for a social media application built with Go.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @tag.name					ops
// @tag.description			Health and operational endpoints
// @tag.name					posts
// @tag.description			Post management
// @tag.name					comments
// @tag.description			Comment management
// @tag.name					users
// @tag.description			User management
// @tag.name					feed
// @tag.description			Feed endpoints
// @BasePath					/v1
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:4000"),
		db: dbConfig{
			addr:          env.GetString("DB_ADDR", "postgres://admin:secret@localhost:5432/e-com?sslmode=disable"),
			maxOpenConns:  env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIddleConns: env.GetInt("DB_MAX_IDDLE_CONNS", 25),
			maxIdleTime:   env.GetDuration("DB_MAX_IDLE_TIME", time.Minute*5),
		},
		env: env.GetString("APP_ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 1,
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
		},
		version: env.GetString("APP_VERSION", "0.0.2"),
		apiURL:  env.GetString("EXTERNAL_URL", "http://localhost:8080"),
	}
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	logger.Infof("starting application in %s mode", cfg.env)
	db, err := db.NewDB(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIddleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Panic(err)
	}
	defer db.Close()
	logger.Infof("db connection pool established")
	store := store.NewStorage(db)

	mailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)
	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
		mailer: mailer,
	}
	mux := app.mount()
	logger.Fatal(app.run(mux))
}
