package pq

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewPqClient(username, password, database string) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", username, password, database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
