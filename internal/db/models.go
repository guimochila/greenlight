// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/guimochila/greenlight/internal/data"
)

type Movie struct {
	ID        uuid.UUID    `db:"id"`
	CreatedAt time.Time    `db:"created_at"`
	Title     string       `db:"title"`
	Year      int32        `db:"year"`
	Runtime   data.Runtime `db:"runtime"`
	Genres    []string     `db:"genres"`
	Version   int32        `db:"version"`
}