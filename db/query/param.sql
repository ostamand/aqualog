-- name: CreateParam :one
INSERT INTO params (
  user_id,
  param_type_id,
  value
) VALUES (
  $1,
  $2,
  $3
) 
RETURNING *;

-- name: GetParam :one
SELECT * FROM params
WHERE id = $1
LIMIT 1;