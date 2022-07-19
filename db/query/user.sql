-- name: CreateUser :one
INSERT INTO users (
  username
) VALUES (
  $1
) 
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 
LIMIT 1;

-- name: GetByUsername :one
SELECT * FROM users
WHERE username = $1 
LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id 
LIMIT $1 
OFFSET $2;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;