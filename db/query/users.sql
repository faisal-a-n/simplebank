-- name: CreateUser :one
INSERT into users (
  "name", "email", "password", "password_changed_at", "created_at"
)
values
($1, $2, $3, $4, $5) RETURNING *;

-- name: GetUser :one
SELECT * from users where id = $1 limit 1;

-- name: ListUsers :many
SELECT * from Users order by id limit $1 offset $2;

-- name: UpdatePassword :one
UPDATE users set password = sqlc.arg(password), password_changed_at = sqlc.arg(passwordChangedAt)
where id = sqlc.arg(id) RETURNING *;