-- name: GetAllUserIdASC :many
SELECT *
FROM users
WHERE user_deleted_at IS NULL
AND (
    sqlc.narg(search)::TEXT IS NULL
    OR sqlc.narg(search) = '::TEXT'
    OR user_email ILIKE '%' || sqlc.narg(search) || '%'
    OR user_name  ILIKE '%' || sqlc.narg(search) || '%'
)
ORDER BY user_id ASC
LIMIT $1 OFFSET $2;

-- name: GetAllUserIdDESC :many
SELECT *
FROM users
WHERE user_deleted_at IS NULL
  AND (
    sqlc.narg(search)::TEXT IS NULL
        OR sqlc.narg(search)::TEXT = ''
        OR user_email ILIKE '%' || sqlc.narg(search) || '%'
    OR user_name  ILIKE '%' || sqlc.narg(search) || '%'
    )
ORDER BY user_id DESC
LIMIT $1 OFFSET $2;

-- name: GetAllUserCreateASC :many
SELECT *
FROM users
WHERE user_deleted_at IS NULL
  AND (
    sqlc.narg(search)::TEXT IS NULL
        OR sqlc.narg(search)::TEXT = ''
        OR user_email ILIKE '%' || sqlc.narg(search) || '%'
    OR user_name  ILIKE '%' || sqlc.narg(search) || '%'
    )
ORDER BY user_created_at ASC
LIMIT $1 OFFSET $2;


-- name: GetAllUserCreateDESC :many
SELECT *
FROM users
WHERE user_deleted_at IS NULL
  AND (
    sqlc.narg(search)::TEXT IS NULL
        OR sqlc.narg(search)::TEXT = ''
        OR user_email ILIKE '%' || sqlc.narg(search) || '%'
    OR user_name  ILIKE '%' || sqlc.narg(search) || '%'
    )
ORDER BY user_created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountAllUsers :one
SELECT COUNT(*) FROM users WHERE user_deleted_at IS NULL;

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
    user_uuid = $1 AND
    user_deleted_at IS NULL
    RETURNING *;

-- name: DeleteUser :execrows
UPDATE users
SET user_deleted_at = NOW()
WHERE user_uuid = $1
AND user_deleted_at IS NULL;

-- name: RestoreUser :execrows
UPDATE users
SET user_deleted_at = NULL
WHERE user_uuid = $1
AND user_deleted_at IS NOT NULL;

-- name: TrashUser :execrows
DELETE
FROM users
WHERE user_uuid = $1
AND user_deleted_at IS NOT NULL;

-- name: GetUserByUUID :one
SELECT * FROM users WHERE user_uuid = $1;