-- name: CreateValue :one
INSERT INTO values (
  user_id,
  value_type_id,
  value
) VALUES (
  $1,
  $2,
  $3
) 
RETURNING *;

-- name: ListValuesPerType :many
SELECT * from values
WHERE value_type_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;