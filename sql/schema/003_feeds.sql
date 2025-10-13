-- +goose Up

CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    url TEXT NOT NULL,
    title TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, url),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE

);

-- +goose Down
DROP TABLE feeds;

