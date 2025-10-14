-- Initial database setup for gomcp
-- This script runs automatically when PostgreSQL container starts

-- Ensure the database is created (docker-compose handles this, but as a safeguard)
CREATE DATABASE IF NOT EXISTS gomcp;

-- Set timezone
SET timezone = 'UTC';

-- Grant necessary permissions
GRANT ALL PRIVILEGES ON DATABASE gomcp TO gomcp;

-- Note: The actual OAuth tables are created automatically by the Go application
-- in internal/storage/storage.go when the app starts
-- This script is just for initial DB setup and extensions if needed

-- Enable extensions if needed in the future
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Log successful initialization
DO $$
BEGIN
  RAISE NOTICE 'Database initialization complete';
END $$;

