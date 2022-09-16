-- name: ListParamOrigins :many
SELECT * 
FROM param_origins
WHERE param_type_name = $1;