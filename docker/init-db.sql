-- Initialize the database with some sample data
-- This file will be executed when the PostgreSQL container starts for the first time

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- The tables will be created automatically by GORM migrations
-- But we can add some initial data here after the application starts

-- Note: This file is mainly for database initialization
-- GORM will handle table creation and migrations automatically
