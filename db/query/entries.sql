-- name: CreateEntry :one
INSERT into entries (
  "account_id", "amount", "created_at"
)
values
($1, $2, $3) RETURNING *;

-- name: GetEntry :one
SELECT * from entries where id = $1 limit 1;

-- name: ListEntries :many
SELECT * from entries order by id limit $1 offset $2;