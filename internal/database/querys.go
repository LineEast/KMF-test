package database

import (
	"context"
)

// CRUD

func (db *Database) Create(ctx context.Context, req, resp []byte) (id int, err error) {
	err = db.Pool.QueryRow(
		ctx,
		"insert into queries (req, resp) values ($1, $2) returning id",
		req, resp,
	).Scan(&id)

	return
}
