-- name: GetAllUsers :many
SELECT * FROM users WHERE user_deleted_at IS NULL;

-- name: CreateUser :one
INSERT INTO users (
    user_email,
    user_password,
    user_name,
    user_age,
    user_status,
    user_role
) VALUES (
     $1, $2, $3, $4, $5, $6
 ) RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
    user_name = COALESCE(sqlc.narg(user_name),user_name),
    user_password = COALESCE(sqlc.narg(user_password),user_password),
    user_age = COALESCE(sqlc.narg(user_age),user_age),
    user_status = COALESCE(sqlc.narg(user_status),user_status),
    user_role = COALESCE(sqlc.narg(user_role),user_role)
WHERE
    user_uuid = sqlc.narg(user_uuid) AND
    user_deleted_at IS NULL
    RETURNING *;