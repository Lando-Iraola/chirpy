-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES(
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING *;

-- name: FindToken :one
SELECT 
    token,
    user_id,
    expires_at,
    revoked_at
FROM
    refresh_tokens
WHERE
    token = $1;

-- name: RevokeRefreshToken :exec
update refresh_tokens 
set 
    revoked_at = NOW(),
    updated_at = NOW()
where 
    token = $1;