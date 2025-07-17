package main

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewHub(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		// Create a mock hub without database for basic functionality tests
		hub := &Hub{
			channels:   make(map[string]*Channel),
			register:   make(chan *Client),
			unregister: make(chan *Client),
			broadcast:  make(chan []byte),
			db:         nil,
			shutdown:   make(chan bool),
		}

		if hub.channels == nil {
			t.Error("Hub channels map should be initialized")
		}
		if hub.register == nil {
			t.Error("Hub register channel should be initialized")
		}
		return
	}
	defer db.Close()

	hub := newHub(db)

	if hub.channels == nil {
		t.Error("Hub channels map should be initialized")
	}
	if hub.register == nil {
		t.Error("Hub register channel should be initialized")
	}
	if hub.unregister == nil {
		t.Error("Hub unregister channel should be initialized")
	}
	if hub.broadcast == nil {
		t.Error("Hub broadcast channel should be initialized")
	}
	if hub.db != db {
		t.Error("Hub database should be set to provided database")
	}
}

func TestHubClientRegistration(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Skip("Skipping hub tests without database")
		return
	}
	defer db.Close()

	hub := newHub(db)

	// Create a mock client
	client := &Client{
		hub:     hub,
		conn:    nil, // We won't actually use the connection in this test
		send:    make(chan []byte, 256),
		channel: "general",
	}

	// Test registration in a separate goroutine
	go func() {
		hub.register <- client
	}()

	// Start hub in background
	go hub.run()
	defer hub.stop() // Ensure cleanup

	// Give some time for registration to process
	time.Sleep(50 * time.Millisecond)

	// Check if client was registered to general channel
	hub.channelsMu.RLock()
	generalChannel, exists := hub.channels["general"]
	hub.channelsMu.RUnlock()

	if !exists {
		t.Error("General channel should exist after client registration")
		return
	}

	generalChannel.clientsMu.RLock()
	_, clientExists := generalChannel.clients[client]
	generalChannel.clientsMu.RUnlock()

	if !clientExists {
		t.Error("Client should be registered in general channel")
	}

	// Test unregistration
	go func() {
		hub.unregister <- client
	}()

	time.Sleep(50 * time.Millisecond)

	// Channel should still exist but client should be gone
	generalChannel.clientsMu.RLock()
	_, clientExists = generalChannel.clients[client]
	generalChannel.clientsMu.RUnlock()

	if clientExists {
		t.Error("Client should be unregistered from general channel")
	}
}

func TestHubBroadcast(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Skip("Skipping hub tests without database")
		return
	}
	defer db.Close()

	hub := newHub(db)

	// Create mock clients
	client1 := &Client{
		hub:     hub,
		conn:    nil,
		send:    make(chan []byte, 256),
		channel: "general",
	}
	client2 := &Client{
		hub:     hub,
		conn:    nil,
		send:    make(chan []byte, 256),
		channel: "general",
	}

	// Register clients
	go hub.run()
	defer hub.stop() // Ensure cleanup

	hub.register <- client1
	hub.register <- client2

	time.Sleep(50 * time.Millisecond)

	// Drain all initial messages (active_channels and join messages)
	// Client1 gets: active_channels, join (for client1), join (for client2)
	// Client2 gets: active_channels, join (for client2)
	for {
		select {
		case <-client1.send:
			// Received message
		case <-time.After(50 * time.Millisecond):
			// No more messages to drain
			goto drainClient2
		}
	}
drainClient2:
	for {
		select {
		case <-client2.send:
			// Received message
		case <-time.After(50 * time.Millisecond):
			// No more messages to drain
			goto testBroadcast
		}
	}
testBroadcast:

	// Test broadcast message
	testMessage := []byte(`{"type":"message","content":"test broadcast"}`)

	go func() {
		hub.broadcast <- testMessage
	}()

	// Check if both clients received the message
	select {
	case msg := <-client1.send:
		if string(msg) != string(testMessage) {
			t.Errorf("Client1 received wrong message: %s", string(msg))
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Client1 did not receive broadcast message")
	}

	select {
	case msg := <-client2.send:
		if string(msg) != string(testMessage) {
			t.Errorf("Client2 received wrong message: %s", string(msg))
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Client2 did not receive broadcast message")
	}
}

func TestSendActiveChannels(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Skip("Skipping hub tests without database")
		return
	}
	defer db.Close()

	hub := newHub(db)

	// Create a persistent channel in database
	err := hub.createChannelInDB("test-persistent", Persistent)
	if err != nil {
		t.Fatalf("Failed to create persistent channel: %v", err)
	}

	// Create an ephemeral channel in memory
	hub.channelsMu.Lock()
	hub.channels["test-ephemeral"] = newChannel("test-ephemeral", Ephemeral)
	hub.channelsMu.Unlock()

	// Create a mock client
	client := &Client{
		hub:     hub,
		conn:    nil,
		send:    make(chan []byte, 256),
		channel: "general",
	}

	// Send active channels
	hub.sendActiveChannels(client)

	// Check if client received active channels message
	select {
	case msg := <-client.send:
		var activeChannelsMsg struct {
			Type     string `json:"type"`
			Channels []struct {
				Name string `json:"name"`
				Type string `json:"type"`
			} `json:"channels"`
		}

		err := json.Unmarshal(msg, &activeChannelsMsg)
		if err != nil {
			t.Fatalf("Failed to unmarshal active channels message: %v", err)
		}

		if activeChannelsMsg.Type != "active_channels" {
			t.Errorf("Expected message type 'active_channels', got '%s'", activeChannelsMsg.Type)
		}

		// Should have at least the persistent channel
		found := false
		for _, ch := range activeChannelsMsg.Channels {
			if ch.Name == "test-persistent" && ch.Type == "persistent" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Active channels should include the persistent test channel")
		}

	case <-time.After(100 * time.Millisecond):
		t.Error("Client did not receive active channels message")
	}
}

func TestHubChannelCleanup(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Skip("Skipping hub tests without database")
		return
	}
	defer db.Close()

	hub := newHub(db)
	go hub.run()
	defer hub.stop() // Ensure cleanup

	// Create an ephemeral channel with a client
	client := &Client{
		hub:     hub,
		conn:    nil,
		send:    make(chan []byte, 256),
		channel: "test-ephemeral",
	}

	hub.register <- client
	time.Sleep(50 * time.Millisecond)

	// Verify channel exists
	hub.channelsMu.RLock()
	_, exists := hub.channels["test-ephemeral"]
	hub.channelsMu.RUnlock()

	if !exists {
		t.Error("Ephemeral channel should exist after client registration")
		return
	}

	// Unregister client
	hub.unregister <- client
	time.Sleep(50 * time.Millisecond)

	// Ephemeral channel should be removed (except general)
	hub.channelsMu.RLock()
	_, exists = hub.channels["test-ephemeral"]
	hub.channelsMu.RUnlock()

	if exists {
		t.Error("Ephemeral channel should be removed after last client disconnects")
	}
}

func TestHubPersistentChannelMemoryCleanup(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Skip("Skipping hub tests without database")
		return
	}
	defer db.Close()

	hub := newHub(db)
	go hub.run()
	defer hub.stop() // Ensure cleanup

	// Create a persistent channel in database
	err := hub.createChannelInDB("test-persistent", Persistent)
	if err != nil {
		t.Fatalf("Failed to create persistent channel: %v", err)
	}

	// Create a client for the persistent channel
	client := &Client{
		hub:     hub,
		conn:    nil,
		send:    make(chan []byte, 256),
		channel: "test-persistent",
	}

	hub.register <- client
	time.Sleep(50 * time.Millisecond)

	// Verify channel exists in memory
	hub.channelsMu.RLock()
	_, exists := hub.channels["test-persistent"]
	hub.channelsMu.RUnlock()

	if !exists {
		t.Error("Persistent channel should exist in memory after client registration")
		return
	}

	// Unregister client
	hub.unregister <- client
	time.Sleep(50 * time.Millisecond)

	// Persistent channel should be removed from memory but preserved in database
	hub.channelsMu.RLock()
	_, exists = hub.channels["test-persistent"]
	hub.channelsMu.RUnlock()

	if exists {
		t.Error("Persistent channel should be removed from memory after last client disconnects")
	}

	// Verify it still exists in database
	var channelType string
	err = db.QueryRow("SELECT type FROM channels WHERE name = $1", "test-persistent").Scan(&channelType)
	if err != nil {
		t.Error("Persistent channel should still exist in database after memory cleanup")
	}
}
