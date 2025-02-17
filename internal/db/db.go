// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.createMovieStmt, err = db.PrepareContext(ctx, createMovie); err != nil {
		return nil, fmt.Errorf("error preparing query CreateMovie: %w", err)
	}
	if q.deleteMovieStmt, err = db.PrepareContext(ctx, deleteMovie); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteMovie: %w", err)
	}
	if q.getAllStmt, err = db.PrepareContext(ctx, getAll); err != nil {
		return nil, fmt.Errorf("error preparing query GetAll: %w", err)
	}
	if q.getMovieStmt, err = db.PrepareContext(ctx, getMovie); err != nil {
		return nil, fmt.Errorf("error preparing query GetMovie: %w", err)
	}
	if q.updateMovieStmt, err = db.PrepareContext(ctx, updateMovie); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateMovie: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createMovieStmt != nil {
		if cerr := q.createMovieStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createMovieStmt: %w", cerr)
		}
	}
	if q.deleteMovieStmt != nil {
		if cerr := q.deleteMovieStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteMovieStmt: %w", cerr)
		}
	}
	if q.getAllStmt != nil {
		if cerr := q.getAllStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllStmt: %w", cerr)
		}
	}
	if q.getMovieStmt != nil {
		if cerr := q.getMovieStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getMovieStmt: %w", cerr)
		}
	}
	if q.updateMovieStmt != nil {
		if cerr := q.updateMovieStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateMovieStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db              DBTX
	tx              *sql.Tx
	createMovieStmt *sql.Stmt
	deleteMovieStmt *sql.Stmt
	getAllStmt      *sql.Stmt
	getMovieStmt    *sql.Stmt
	updateMovieStmt *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:              tx,
		tx:              tx,
		createMovieStmt: q.createMovieStmt,
		deleteMovieStmt: q.deleteMovieStmt,
		getAllStmt:      q.getAllStmt,
		getMovieStmt:    q.getMovieStmt,
		updateMovieStmt: q.updateMovieStmt,
	}
}
