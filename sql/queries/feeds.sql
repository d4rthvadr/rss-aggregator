-- name: CreateFeed :one
INSERT INTO feeds (id, title, url, user_id) 
VALUES ($1, $2, $3, $4) 
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds ORDER BY created_at DESC;