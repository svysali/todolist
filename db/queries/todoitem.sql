-- name: CreateItem :one
INSERT INTO items (
  title
) VALUES (
  $1
) RETURNING *;

-- name: GetItem :one
SELECT * FROM items
WHERE id = $1;

-- name: ListItems :many
SELECT * FROM items
ORDER BY created_at;
