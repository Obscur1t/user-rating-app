-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    nickname TEXT NOT NULL UNIQUE,
    likes INT NOT NULL DEFAULT 0 CHECK (likes >= 0),
    viewers INT NOT NULL DEFAULT 0 CHECK (viewers >= 0),

    CONSTRAINT likes_lte_viewers CHECK (likes <= viewers),

    rating NUMERIC GENERATED ALWAYS AS (
        CASE WHEN viewers > 0
             THEN ROUND(likes::NUMERIC / viewers, 3)
             ELSE 0
        END
    ) STORED
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
