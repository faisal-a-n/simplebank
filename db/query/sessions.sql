-- name: CreateSession :one
INSERT into sessions (
  "id", "user_id", "refresh_token", "user_agent", "client_ip", "expires_at", "created_at"
)
values
($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetSession :one
SELECT * from sessions where id = $1 LIMIT 1;

-- name: UpdateSession :exec
UPDATE sessions set is_blocked = $1 where user_id = $2;