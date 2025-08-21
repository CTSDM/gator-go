-- name: CreateFeed :one
INSERT INTO feeds (created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFeeds :many
SELECT feeds.name, feeds.url, u.name as username
FROM feeds
INNER JOIN users as u
ON u.id = feeds.user_id;
