-- +goose Up
DROP TABLE IF EXISTS otps;
DROP INDEX IF EXISTS idx_otps_user_type;

-- +goose Down

CREATE TABLE otps (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  hashed_otp TEXT NOT NULL,
  type otp_type NOT NULL,
  is_used BOOLEAN DEFAULT FALSE,
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_otps_user_type ON otps(user_id, type);