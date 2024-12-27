// Copyleft (c) 2024, guimochila. Happy Coding.
package main

import (
	"database/sql"
	"fmt"
	"net/http"

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

	movie := db.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()

	if validator.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)

		return
	}

	_, err = app.querier.CreateMovie(r.Context(), db.CreateMovieParams{
		Title:   movie.Title,
		Year:    movie.Year,
		Runtime: movie.Runtime,
		Genres:  movie.Genres,
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)

		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)

		return
	}

	movie, err := app.querier.GetMovie(r.Context(), id)
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
	if validator.ValidateMovie(v, params); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)

		return
	}

	movie, err := app.querier.UpdateMovie(r.Context(), params)
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
