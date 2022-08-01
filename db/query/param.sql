-- name: CreateParam :one
INSERT INTO params (
  user_id,
  param_type_id,
  value,
  timestamp
) VALUES (
  $1,
  $2,
  $3,
  $4
) 
RETURNING *;

-- name: GetParam :one
SELECT * FROM params
WHERE id = $1
LIMIT 1;

-- name: GetParams :many
SELECT p.id, p.user_id, p.value, p.timestamp, p.created_at 
FROM params as p
INNER JOIN param_types as t ON p.param_type_id = t.id
WHERE p.user_id = sqlc.arg(user_id)
	AND t.name = sqlc.arg(param_type_name)
	AND p.timestamp >= sqlc.arg('from')
	AND p.timestamp < sqlc.arg('to')
ORDER BY p.timestamp
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');