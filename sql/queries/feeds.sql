-- name: CreateFeed :one
INSERT INTO feeds (created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFeeds :many
SELECT feeds.name, feeds.url, u.name as username
FROM feeds
INNER JOIN users as u
ON u.id = feeds.user_id;

-- name: GetFeedByURL :one
SELECT *
FROM feeds
WHERE url = $1;

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = $1,
    updated_at = $2
WHERE id = $3
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
