-- +goose Up
CREATE TYPE otp_type AS ENUM (
  'account_verification',
  'forgot_password'
);

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

-- +goose Down
DROP INDEX IF EXISTS idx_otps_user_type;
DROP TABLE IF EXISTS otps;
DROP TYPE IF EXISTS otp_type;