// Copyleft (c) 2024, guimochila. Happy Coding.
package querier

import (
	"database/sql"

	"github.com/guimochila/greenlight/internal/db"
)

type Querier struct {
	*db.Queries
	db *sql.DB
}

func New(conn *sql.DB) *Querier {
	return &Querier{
		Queries: db.New(conn),
		db:      conn,
	}
}
