-- name: CreateAccount :one
INSERT into accounts (
  "user_id", "name", "balance", "currency", "created_at"
)
values
($1, $2, $3, $4, $5) RETURNING *;

-- name: GetAccount :one
SELECT * from accounts where id = $1 limit 1;

-- name: GetAccountForUpdate :one
SELECT * from accounts where id = $1 limit 1 for NO KEY UPDATE;

-- name: ListAccounts :many
SELECT * from accounts order by id limit $1 offset $2;

-- name: ListAccountsForUser :many
SELECT * from accounts where user_id = $1 order by id limit $2 offset $3;

-- name: UpdateBalance :one
UPDATE accounts set balance = balance + sqlc.arg(amount) where id = sqlc.arg(id) RETURNING *;

-- name: DeleteAccount :exec
DELETE from accounts where id = $1;