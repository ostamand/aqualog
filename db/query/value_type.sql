-- name: GetValueTypeByName :one
SELECT * FROM value_types
WHERE user_id = $1 AND name = $2
LIMIT 1;

-- name: GetValueType :one
SELECT * FROM value_types
WHERE id = $1
LIMIT 1;

-- name: CreateValueType :one
INSERT INTO value_types (
  name,
  description,
  unit,
  user_id,
  target,
  min,
  max
) VALUES (
  $1, 
  $2, 
  $3,
  $4,
  $5,
  $6,
  $7
) 
RETURNING *;