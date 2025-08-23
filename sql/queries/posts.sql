-- name: CreatePost :one
INSERT INTO posts (created_at, updated_at, title, url, description, published_at, feed_id)
VALUES($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT DO NOTHING
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.*
FROM posts
INNER JOIN feeds
ON feeds.id = posts.feed_id
INNER JOIN feed_follows
ON feed_follows.user_id = $1
ORDER BY posts.updated_at DESC NULLS LAST
LIMIT $2;
