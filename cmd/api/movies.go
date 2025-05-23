// Copyleft (c) 2024, guimochila. Happy Coding.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/guimochila/greenlight/config"
	"github.com/guimochila/greenlight/internal/data"
	"github.com/guimochila/greenlight/internal/db"
	"github.com/guimochila/greenlight/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)

		return
	}

	params := db.CreateMovieParams{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()

	if data.ValidateMovie(v, params); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)

		return
	}

	ctx, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), config.MaxQueryTimeout)
	defer cancel()

	movie, err := app.querier.CreateMovie(ctx, params)
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%s", movie.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)

		return
	}

	ctx, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), config.MaxQueryTimeout)
	defer cancel()

	movie, err := app.querier.GetMovie(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)

		return
	}

	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)

		return
	}

	params := db.UpdateMovieParams{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
		ID:      id,
	}

	v := validator.New()
	if data.ValidateMovie(v, params); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)

		return
	}

	ctx, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), config.MaxQueryTimeout)
	defer cancel()

	movie, err := app.querier.UpdateMovie(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)

		return
	}

	ctx, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), config.MaxQueryTimeout)
	defer cancel()

	err = app.querier.DeleteMovie(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			app.notFoundResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "created_at")
	input.Filters.SortSafeList = []string{
		"created_at",
		"title",
		"year",
		"runtime",
		"-created_at",
		"-title",
		"-year",
		"-runtime",
	}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)

		return
	}

	ctx, cancel := context.WithTimeout(context.WithoutCancel(r.Context()), config.MaxQueryTimeout)
	defer cancel()

	dbMovies, err := app.querier.GetAll(ctx, db.GetAllParams{
		PlaintoTsquery: input.Title,
		Genres:         input.Genres,
		SortColumn: sql.NullString{
			String: input.Filters.Sort,
			Valid:  true,
		},
		Limit:  input.Filters.Limit(),
		Offset: input.Filters.Offset(),
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}

	var totalCount int64
	movies := make([]db.GetAllRow, 0, len(dbMovies))

	if len(dbMovies) > 0 {
		totalCount = dbMovies[0].TotalCount
	}

	for _, m := range dbMovies {
		movies = append(movies, m)
	}

	metadata := data.CalculateMetadata(int(totalCount), input.Page, input.PageSize)

	err = app.writeJSON(w, http.StatusOK, envelope{"movies": movies, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
