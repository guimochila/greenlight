-- name: CreateMovie :one
INSERT INTO movies (title, year, runtime, genres)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, version;

-- name: GetMovie :one
SELECT id, created_at, title, year, runtime, genres, version
FROM movies
WHERE id = $1;

-- name: GetAll :many
SELECT id, created_at, title, year, runtime, genres, version
FROM movies
WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
AND (genres @> $2 OR $2 = '{}')
ORDER BY
  CASE WHEN sqlc.narg('sort_column')::text = 'title' THEN title END ASC,
  CASE WHEN sqlc.narg('sort_column')::text = '-title' THEN title END DESC,
  CASE WHEN sqlc.narg('sort_column')::text = 'runtime' THEN runtime END ASC,
  CASE WHEN sqlc.narg('sort_column')::text = '-runtime' THEN runtime END DESC,
  CASE WHEN sqlc.narg('sort_column')::text = 'year' THEN year END ASC,
  CASE WHEN sqlc.narg('sort_column')::text = '-year' THEN year END DESC,
  CASE WHEN sqlc.narg('sort_column')::text = 'created_at' THEN created_at END ASC,
  CASE WHEN sqlc.narg('sort_column')::text = '-created_at' THEN created_at END DESC,
id ASC;

-- name: UpdateMovie :one
UPDATE movies
SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
WHERE id = $5 AND version = $6
RETURNING id, created_at, year, runtime, genres, version;

-- name: DeleteMovie :exec
DELETE FROM movies
WHERE id = $1;
