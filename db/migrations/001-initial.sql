-- +migrate Up

-- Create 'items' table
CREATE TABLE items (
  id SERIAL PRIMARY KEY,
  title VARCHAR  NOT NULL,
  created_at timestamptz NOT NULL DEFAULT (now()),
  completed BOOLEAN NOT NULL DEFAULT FALSE
);

-- +migrate Down
DROP TABLE IF EXISTS items;
