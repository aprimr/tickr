-- +goose Up
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    hashed_token TEXT NOT NULL,
    device_info TEXT,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id
on refresh_tokens(user_id);

CREATE INDEX idx_refresh_tokens_hashed_token
ON refresh_tokens(hashed_token);

-- +goose Down
DROP INDEX IF EXISTS idx_refresh_tokens_hashed_token;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP TABLE IF EXISTS refresh_tokens;