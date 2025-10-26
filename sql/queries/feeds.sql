-- name: CreateFeed :one
INSERT INTO feeds (id, title, url, user_id) 
VALUES ($1, $2, $3, $4) 
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds ORDER BY created_at DESC;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds 
ORDER BY last_fetched_at ASC NULLS FIRST 
LIMIT $1;

-- name: UpdateFeedLastFetchedAt :exec
UPDATE feeds 
SET last_fetched_at = NOW(), updated_at = CURRENT_TIMESTAMP 
WHERE id = $1
RETURNING *;