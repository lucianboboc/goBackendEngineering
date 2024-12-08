package main

import (
	"goBackendEngineering/internal/db"
	"goBackendEngineering/internal/env"
	store2 "goBackendEngineering/internal/store"
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
	store := store2.NewPostgresStorage(conn)
	db.Seed(store)
}
