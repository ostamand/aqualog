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

-- name: GetParamByID :one
SELECT 
p.id as param_id,
t.id as param_type_id,
p."value",
p.timestamp,
t."name",
CASE WHEN t.target IS NULL THEN -999 ELSE t.target END AS target,
CASE WHEN t."min" IS NULL THEN -999 ELSE t."min" END AS "min",
CASE WHEN t."max" IS NULL THEN -999 ELSE t."max" END AS "max",
p.created_at
FROM params as p
INNER JOIN param_types AS t ON p.param_type_id = t.id
WHERE p.user_id=$2 AND p.id = $1
LIMIT 1;

-- name: ListParamsByType :many
SELECT p.id, p.value, p.timestamp, p.created_at 
FROM params as p
INNER JOIN param_types as t ON p.param_type_id = t.id
WHERE p.user_id = sqlc.arg(user_id)
	AND t.name = sqlc.arg(param_type_name)
	AND p.timestamp >= sqlc.arg('from')
	AND p.timestamp < sqlc.arg('to')
ORDER BY p.timestamp DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');