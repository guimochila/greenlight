package data

import (
	"reflect"
	"time"

	"github.com/guimochila/greenlight/internal/validator"
)

const (
	TitleBytesLong = 500
	MinimunYear    = 1888
	MaxGenres      = 5
)

func ValidateMovie(v *validator.Validator, movie any) {
	val := reflect.ValueOf(movie)

	if val.Kind() != reflect.Struct {
		v.Check(false, "movie", "not a struct")

		return
	}

	if title := val.FieldByName("Title"); title.IsValid() && title.Kind() == reflect.String {
		titleStr := title.String()
		v.Check(titleStr != "", "title", "must be provided")
		v.Check(len(titleStr) <= TitleBytesLong, "title", "must not be more than 500 bytes long")
	}

	if year := val.FieldByName("Year"); year.IsValid() && year.Kind() == reflect.Int32 {
		yearInt := int32(year.Int())
		v.Check(yearInt != 0, "year", "must be provided")
		v.Check(yearInt >= MinimunYear, "year", "must be greater than 1888")
		//nolint:all
		v.Check(yearInt <= int32(time.Now().Year()), "year", "must not be in the future")
	}

	if runtime := val.FieldByName("Runtime"); runtime.IsValid() && runtime.Kind() == reflect.Int32 {
		runtimeInt := Runtime(runtime.Int())
		v.Check(runtimeInt != 0, "runtime", "must be provided")
		v.Check(runtimeInt > 0, "runtime", "must be a positive integer")
	}

	if genres := val.FieldByName("Genres"); genres.IsValid() && genres.Kind() == reflect.Slice {
		genresSlice := reflect.ValueOf(genres.Interface())
		v.Check(genresSlice.Len() >= 1, "genres", "must contain at least 1 genre")
		v.Check(genresSlice.Len() <= MaxGenres, "genres", "must not contain more than 5 genres")

		uniqueGenres := make([]string, genresSlice.Len())
		for i := 0; i < genresSlice.Len(); i++ {
			uniqueGenres[i] = genresSlice.Index(i).String()
		}

		v.Check(validator.Unique(uniqueGenres), "genres", "must not contain duplicate values")
	}
}
