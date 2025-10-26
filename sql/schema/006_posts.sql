-- +goose Up

CREATE TABLE posts (
    id UUID PRIMARY KEY,
    url TEXT NOT NULL,
    userId UUID REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    published_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(userId, feed_id),
    feed_id UUID REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;

