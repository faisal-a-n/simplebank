-- name: CreateTransaction :one
INSERT into transactions (
 "from_account_id", "to_account_id", "from_entry_id", "to_entry_id", "amount", "created_at"
)
values
($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetTransaction :one
SELECT * from transactions where id = $1 limit 1;

-- name: ListTransactions :many
SELECT * from transactions order by id limit $1 offset $2;