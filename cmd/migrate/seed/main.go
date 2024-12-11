package main

import (
	"github.com/lucianboboc/goBackendEngineering/internal/db"
	"github.com/lucianboboc/goBackendEngineering/internal/env"
	"github.com/lucianboboc/goBackendEngineering/internal/store"
)

func main() {
	conn, err := db.New(
		env.GetString("DB_DSN", "postgres://postgres:postgres@localhost/postgres?sslmode=disable"),
		3,
		3,
		"15m",
	)
	if err != nil {
		panic(err)
	}
	storage := store.NewPostgresStorage(conn)
	db.Seed(storage, conn)
}
