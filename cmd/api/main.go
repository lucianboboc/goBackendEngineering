package main

import (
	"github.com/joho/godotenv"
	"goBackendEngineering/internal/db"
	"goBackendEngineering/internal/env"
	"goBackendEngineering/internal/store"
	"log"
	"log/slog"
	"os"
)

const version = "0.0.1"

type Post struct {
	Title string `json:"title" validate:"required,max=100"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr: env.GetString("ADDR", "8080"),
		db: dbConfig{
			dsn:          env.GetString("DB_DSN", "postgres://postgres:postgres@localhost/postgres?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConn:  env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "development"),
	}

	db, err := db.New(
		cfg.db.dsn,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConn,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}

	defer db.Close()
	log.Println("database connection established...")
	store := store.NewPostgresStorage(db)

	app := &application{
		config: cfg,
		store:  store,
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	mux := app.mount()
	if err := app.run(mux); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
