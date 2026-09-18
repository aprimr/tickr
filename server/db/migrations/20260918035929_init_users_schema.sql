-- +goose Up

-- Enums
CREATE TYPE user_role AS ENUM ('user', 'venue_admin', 'super_admin');
CREATE TYPE venue_status AS ENUM ('pending', 'approved', 'rejected');

-- Core users table
CREATE TABLE users (
  id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  email TEXT UNIQUE NOT NULL,
  phone_number TEXT NOT NULL,
  password_hash TEXT NOT NULL, 
  role user_role NOT NULL DEFAULT 'user',
  is_active BOOLEAN DEFAULT TRUE,
  is_email_verified BOOLEAN DEFAULT FALSE,
  is_phone_verified BOOLEAN DEFAULT FALSE,
  last_login_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Users detail table
CREATE TABLE user_details (
  id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  full_name TEXT NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Venue detail table
CREATE TABLE venue_details (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    venue_name TEXT NOT NULL,
    address TEXT NOT NULL,
    city TEXT NOT NULL,
    total_screens INTEGER NOT NULL DEFAULT 1,
    venue_status venue_status NOT NULL DEFAULT 'pending',
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    approved_at TIMESTAMPTZ NULL,
    rejected_by UUID REFERENCES users(id) ON DELETE SET NULL,
    rejected_at TIMESTAMPTZ NULL,
    venue_created_at TIMESTAMPTZ DEFAULT NOW(),
    venue_updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Platform admin detais table
CREATE TABLE admin_details (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    full_name TEXT NOT NULL,
    is_root_user BOOLEAN DEFAULT FALSE,
    admin_created_at TIMESTAMPTZ DEFAULT NOW(),
    admin_updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose Down

-- Tables
DROP TABLE IF EXISTS admin_details;
DROP TABLE IF EXISTS venue_details;
DROP TABLE IF EXISTS user_details;
DROP TABLE IF EXISTS user;

-- Enums
DROP TYPE IF EXISTS venue_status;
DROP TYPE IF EXISTS user_role;