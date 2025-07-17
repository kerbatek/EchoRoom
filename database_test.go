package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Load .env file for test configuration
	if err := godotenv.Load(); err != nil {
		// Silent fail - not required for tests
	}
	
	// Try TEST_DATABASE_URL first
	testConnStr := os.Getenv("TEST_DATABASE_URL")
	if testConnStr == "" {
		// Build from individual environment variables
		host := os.Getenv("TEST_DB_HOST")
		port := os.Getenv("TEST_DB_PORT")
		user := os.Getenv("TEST_DB_USER")
		password := os.Getenv("TEST_DB_PASSWORD")
		dbname := os.Getenv("TEST_DB_NAME")
		sslmode := os.Getenv("TEST_DB_SSLMODE")
		
		// Use defaults if not set
		if host == "" {
			host = "localhost"
		}
		if port == "" {
			port = "5432"
		}
		if user == "" {
			user = "postgres"
		}
		if password == "" {
			password = "password"
		}
		if dbname == "" {
			dbname = "chat_app_test"
		}
		if sslmode == "" {
			sslmode = "disable"
		}
		
		testConnStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&connect_timeout=3",
			user, password, host, port, dbname, sslmode)
	}

	db, err := sql.Open("postgres", testConnStr)
	if err != nil {
		t.Skipf("Skipping database tests - PostgreSQL not available: %v", err)
		return nil
	}

	// Test connection
	if err := db.Ping(); err != nil {
		t.Skipf("Skipping database tests - PostgreSQL not reachable: %v", err)
		return nil
	}

	// Clean up tables if they exist
	db.Exec("DROP TABLE IF EXISTS messages")
	db.Exec("DROP TABLE IF EXISTS channels")

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS channels (
			name VARCHAR(100) PRIMARY KEY,
			type VARCHAR(20) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create channels table: %v", err)
	}

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
		t.Fatalf("Failed to create messages table: %v", err)
	}

	return db
}

func TestInitDatabase(t *testing.T) {
	db, err := initDatabase()
	if err != nil {
		t.Skipf("Skipping database initialization test: %v", err)
		return
	}
	defer db.Close()

	// Test that tables exist
	var tableName string
	err = db.QueryRow("SELECT table_name FROM information_schema.tables WHERE table_name = 'channels'").Scan(&tableName)
	if err != nil {
		t.Errorf("Channels table not found: %v", err)
	}

	err = db.QueryRow("SELECT table_name FROM information_schema.tables WHERE table_name = 'messages'").Scan(&tableName)
	if err != nil {
		t.Errorf("Messages table not found: %v", err)
	}
}

func TestSaveMessage(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	hub := newHub(db)

	// Create a persistent channel first
	err := hub.createChannelInDB("test-persistent", Persistent)
	if err != nil {
		t.Fatalf("Failed to create persistent channel: %v", err)
	}

	// Test saving message to persistent channel
	msg := Message{
		Username:  "testuser",
		Content:   "test message",
		Type:      "message",
		Channel:   "test-persistent",
		Timestamp: time.Now().UTC(),
	}

	err = hub.saveMessage(msg)
	if err != nil {
		t.Errorf("Failed to save message: %v", err)
	}

	// Test that message was saved
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM messages WHERE channel_name = $1", "test-persistent").Scan(&count)
	if err != nil {
		t.Errorf("Failed to query messages: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 message, got %d", count)
	}

	// Test that ephemeral channel messages are not saved
	msg.Channel = "ephemeral-channel"
	err = hub.saveMessage(msg)
	if err == nil {
		t.Errorf("Expected error when saving to non-existent channel")
	}
}

func TestGetChannelHistory(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	hub := newHub(db)

	// Create a persistent channel
	err := hub.createChannelInDB("history-test", Persistent)
	if err != nil {
		t.Fatalf("Failed to create persistent channel: %v", err)
	}

	// Add some test messages
	messages := []Message{
		{Username: "user1", Content: "message 1", Type: "message", Channel: "history-test"},
		{Username: "user2", Content: "message 2", Type: "message", Channel: "history-test"},
		{Username: "user3", Content: "message 3", Type: "message", Channel: "history-test"},
	}

	for _, msg := range messages {
		err = hub.saveMessage(msg)
		if err != nil {
			t.Fatalf("Failed to save message: %v", err)
		}
		time.Sleep(10 * time.Millisecond) // Ensure different timestamps
	}

	// Test getting history
	history, err := hub.getChannelHistory("history-test", 10)
	if err != nil {
		t.Errorf("Failed to get channel history: %v", err)
	}

	if len(history) != 3 {
		t.Errorf("Expected 3 messages in history, got %d", len(history))
	}

	// Test that messages are in reverse chronological order (newest first)
	if history[0].Content != "message 3" { // Should be newest first
		t.Errorf("Expected first message to be 'message 3', got '%s'", history[0].Content)
	}
}

func TestGetChannelType(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	hub := newHub(db)

	// Create channels of different types
	err := hub.createChannelInDB("persistent-test", Persistent)
	if err != nil {
		t.Fatalf("Failed to create persistent channel: %v", err)
	}

	// Test getting persistent channel type
	channelType, err := hub.getChannelType("persistent-test")
	if err != nil {
		t.Errorf("Failed to get channel type: %v", err)
	}
	if channelType != Persistent {
		t.Errorf("Expected Persistent, got %s", channelType)
	}

	// Test getting non-existent channel type
	_, err = hub.getChannelType("non-existent")
	if err == nil {
		t.Errorf("Expected error for non-existent channel")
	}
}

func TestCreateChannelInDB(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	hub := newHub(db)

	// Test creating persistent channel
	err := hub.createChannelInDB("new-persistent", Persistent)
	if err != nil {
		t.Errorf("Failed to create persistent channel: %v", err)
	}

	// Verify channel was created
	var channelType string
	err = db.QueryRow("SELECT type FROM channels WHERE name = $1", "new-persistent").Scan(&channelType)
	if err != nil {
		t.Errorf("Failed to query created channel: %v", err)
	}
	if channelType != string(Persistent) {
		t.Errorf("Expected 'persistent', got '%s'", channelType)
	}

	// Test creating ephemeral channel (should not be stored)
	err = hub.createChannelInDB("new-ephemeral", Ephemeral)
	if err != nil {
		t.Errorf("Unexpected error creating ephemeral channel: %v", err)
	}

	// Verify ephemeral channel was NOT created in database
	err = db.QueryRow("SELECT type FROM channels WHERE name = $1", "new-ephemeral").Scan(&channelType)
	if err == nil {
		t.Errorf("Ephemeral channel should not be stored in database")
	}
}

func TestGetDefaultDBConfig(t *testing.T) {
	// Save original environment
	originalHost := os.Getenv("DB_HOST")
	originalPort := os.Getenv("DB_PORT")

	// Test with custom environment variables
	os.Setenv("DB_HOST", "custom-host")
	os.Setenv("DB_PORT", "5433")

	config := getDefaultDBConfig()

	if config.Host != "custom-host" {
		t.Errorf("Expected Host 'custom-host', got '%s'", config.Host)
	}
	if config.Port != "5433" {
		t.Errorf("Expected Port '5433', got '%s'", config.Port)
	}

	// Test defaults
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")

	config = getDefaultDBConfig()

	if config.Host != "localhost" {
		t.Errorf("Expected default Host 'localhost', got '%s'", config.Host)
	}
	if config.Port != "5432" {
		t.Errorf("Expected default Port '5432', got '%s'", config.Port)
	}

	// Restore original environment
	if originalHost != "" {
		os.Setenv("DB_HOST", originalHost)
	}
	if originalPort != "" {
		os.Setenv("DB_PORT", originalPort)
	}
}