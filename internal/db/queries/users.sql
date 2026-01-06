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
