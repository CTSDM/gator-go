-- name: CreateFeedFollows :one
WITH new_inserted_feed_follows AS (
    INSERT INTO feed_follows (created_at, updated_at, user_id, feed_id)
    VALUES($1, $2, $3, $4)
    RETURNING *
)
SELECT new_inserted_feed_follows.*, u.name as user_name, f.name as feed_name
FROM new_inserted_feed_follows
INNER JOIN users as u
ON u.id = new_inserted_feed_follows.user_id
INNER JOIN feeds as f
ON f.id = new_inserted_feed_follows.feed_id;

-- name: GetFeedfollowsForUser :many
SELECT feed_follows.*, users.name AS user_name, feeds.name AS feed_name
FROM feed_follows
INNER JOIN users
ON users.id = feed_follows.user_id
INNER JOIN feeds
ON feeds.id = feed_follows.feed_id
WHERE users.name = $1;
