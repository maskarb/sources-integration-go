package database

import "github.com/go-pg/pg/v10"

func NewConnection() *pg.DB {
	return pg.Connect(&pg.Options{
		Addr:     "localhost:15432",
		User:     "postgres",
		Password: "postgres",
		Database: "postgres",
	})
}
