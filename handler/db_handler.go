package handler

import "github.com/jmoiron/sqlx"

type dbHandler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *dbHandler {
	return &dbHandler{db: db}
}
