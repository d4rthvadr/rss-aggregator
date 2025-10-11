-- name: CreateUser :one
INSERT INTO users (id, name, created_at) 
VALUES ($1, $2, $3) 
RETURNING *;