// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"
)

type Querier interface {
	CreateMovie(ctx context.Context, arg CreateMovieParams) (CreateMovieRow, error)
}

var _ Querier = (*Queries)(nil)