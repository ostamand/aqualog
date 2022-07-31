-- name: GetParamTypeByName :one
SELECT * FROM param_types
WHERE user_id = $1 AND name = $2
LIMIT 1;

-- name: GetParamType :one
SELECT * FROM param_types
WHERE id = $1
LIMIT 1;

-- name: CreateValueType :one
INSERT INTO param_types (
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