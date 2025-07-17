-- Database initialization script for EchoRoom
-- This script sets up the initial database schema and data

-- Create channels table if it doesn't exist
CREATE TABLE IF NOT EXISTS channels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('ephemeral', 'persistent')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create messages table if it doesn't exist
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    channel VARCHAR(100) NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (channel) REFERENCES channels(name) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_messages_channel ON messages(channel);
CREATE INDEX IF NOT EXISTS idx_messages_timestamp ON messages(timestamp);
CREATE INDEX IF NOT EXISTS idx_channels_name ON channels(name);
CREATE INDEX IF NOT EXISTS idx_channels_type ON channels(type);

-- Insert default channels
INSERT INTO channels (name, type) VALUES ('general', 'persistent') ON CONFLICT (name) DO NOTHING;
INSERT INTO channels (name, type) VALUES ('announcements', 'persistent') ON CONFLICT (name) DO NOTHING;
INSERT INTO channels (name, type) VALUES ('support', 'persistent') ON CONFLICT (name) DO NOTHING;

-- Insert welcome message
INSERT INTO messages (username, content, type, channel) VALUES 
('System', 'Welcome to EchoRoom! 🎉', 'system_message', 'general') 
ON CONFLICT DO NOTHING;

-- Grant necessary permissions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO postgres;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO postgres;