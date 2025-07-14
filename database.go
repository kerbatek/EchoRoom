package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	
	_ "github.com/lib/pq"
)


func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDefaultDBConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "chat_app"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func initDatabase() (*sql.DB, error) {
	config := getDefaultDBConfig()

	// Build PostgreSQL connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Create channels table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS channels (
			name VARCHAR(100) PRIMARY KEY,
			type VARCHAR(20) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create channels table: %v", err)
	}

	// Create messages table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id SERIAL PRIMARY KEY,
			channel_name VARCHAR(100) NOT NULL,
			username VARCHAR(100) NOT NULL,
			content TEXT NOT NULL,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (channel_name) REFERENCES channels (name)
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create messages table: %v", err)
	}

	// Note: general channel is ephemeral and not stored in database

	log.Printf("Successfully connected to PostgreSQL database: %s@%s:%s/%s",
		config.User, config.Host, config.Port, config.DBName)

	return db, nil
}

func (h *Hub) saveMessage(msg Message) error {
	if msg.Type != "message" {
		return nil // Only save regular messages
	}

	// This function should only be called for persistent channels
	_, err := h.db.Exec(`
		INSERT INTO messages (channel_name, username, content, timestamp) 
		VALUES ($1, $2, $3, $4)
	`, msg.Channel, msg.Username, msg.Content, msg.Timestamp)

	return err
}

func (h *Hub) getChannelHistory(channelName string, limit int) ([]Message, error) {
	rows, err := h.db.Query(`
		SELECT id, username, content, timestamp 
		FROM messages 
		WHERE channel_name = $1 
		ORDER BY timestamp DESC 
		LIMIT $2
	`, channelName, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.ID, &msg.Username, &msg.Content, &msg.Timestamp)
		if err != nil {
			return nil, err
		}
		msg.Channel = channelName
		msg.Type = "message"
		messages = append([]Message{msg}, messages...) // Reverse order
	}

	return messages, nil
}

func (h *Hub) getChannelType(channelName string) (ChannelType, error) {
	var channelType string
	err := h.db.QueryRow("SELECT type FROM channels WHERE name = $1", channelName).Scan(&channelType)
	if err != nil {
		return Ephemeral, err
	}
	return ChannelType(channelType), nil
}

func (h *Hub) createChannelInDB(name string, channelType ChannelType) error {
	// Only store persistent channels in database
	if channelType == Persistent {
		_, err := h.db.Exec("INSERT INTO channels (name, type) VALUES ($1, $2) ON CONFLICT (name) DO NOTHING", name, string(channelType))
		return err
	}
	return nil
}