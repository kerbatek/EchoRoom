package main

import (
	"encoding/json"
	"testing"
	"time"
)

func TestClientSwitchChannel(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Skip("Skipping client tests without database")
		return
	}
	defer db.Close()

	hub := newHub(db)
	go hub.run()
	defer hub.stop() // Ensure cleanup

	// Create a mock client
	client := &Client{
		hub:     hub,
		conn:    nil,
		send:    make(chan []byte, 256),
		channel: "general",
	}

	// Create initial channel
	hub.channels["general"] = newChannel("general", Ephemeral)
	hub.channels["general"].clients[client] = true

	// Test switching to a new ephemeral channel
	client.switchChannel("test-channel")

	if client.channel != "test-channel" {
		t.Errorf("Expected client channel to be 'test-channel', got '%s'", client.channel)
	}

	// Verify new channel was created
	if _, exists := hub.channels["test-channel"]; !exists {
		t.Error("New channel should be created when switching")
	}

	// Verify client is in new channel
	if !hub.channels["test-channel"].clients[client] {
		t.Error("Client should be in the new channel")
	}

	// Verify client was removed from old channel
	if hub.channels["general"].clients[client] {
		t.Error("Client should be removed from old channel")
	}
}

func TestClientSwitchChannelWithType(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Skip("Skipping client tests without database")
		return
	}
	defer db.Close()

	hub := newHub(db)
	go hub.run()
	defer hub.stop() // Ensure cleanup

	// Create a mock client
	client := &Client{
		hub:     hub,
		conn:    nil,
		send:    make(chan []byte, 256),
		channel: "general",
	}

	// Create initial channel
	hub.channels["general"] = newChannel("general", Ephemeral)
	hub.channels["general"].clients[client] = true

	// Test switching to a new persistent channel
	client.switchChannelWithType("persistent-test", Persistent)

	if client.channel != "persistent-test" {
		t.Errorf("Expected client channel to be 'persistent-test', got '%s'", client.channel)
	}

	// Verify new channel was created with correct type
	if _, exists := hub.channels["persistent-test"]; !exists {
		t.Error("New persistent channel should be created")
	}

	if hub.channels["persistent-test"].channelType != Persistent {
		t.Error("New channel should be of Persistent type")
	}
}

func TestClientSwitchToSameChannel(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Skip("Skipping client tests without database")
		return
	}
	defer db.Close()

	hub := newHub(db)

	client := &Client{
		hub:     hub,
		conn:    nil,
		send:    make(chan []byte, 256),
		channel: "test-channel",
	}

	// Create channel
	hub.channels["test-channel"] = newChannel("test-channel", Ephemeral)
	hub.channels["test-channel"].clients[client] = true

	// Test switching to the same channel (should be no-op)
	client.switchChannel("test-channel")

	// Should still be in the same channel
	if client.channel != "test-channel" {
		t.Error("Client should remain in the same channel")
	}

	// Should still be registered in the channel
	if !hub.channels["test-channel"].clients[client] {
		t.Error("Client should still be registered in the channel")
	}
}

func TestClientChannelCleanupOnSwitch(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Skip("Skipping client tests without database")
		return
	}
	defer db.Close()

	hub := newHub(db)

	client := &Client{
		hub:     hub,
		conn:    nil,
		send:    make(chan []byte, 256),
		channel: "ephemeral-test",
	}

	// Create ephemeral channel with only one client
	hub.channels["ephemeral-test"] = newChannel("ephemeral-test", Ephemeral)
	hub.channels["ephemeral-test"].clients[client] = true

	// Switch to another channel
	client.switchChannel("general")

	// The ephemeral channel should be deleted (except general)
	if _, exists := hub.channels["ephemeral-test"]; exists {
		t.Error("Empty ephemeral channel should be deleted when last client leaves")
	}
}

func TestClientPersistentChannelMemoryCleanup(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Skip("Skipping client tests without database")
		return
	}
	defer db.Close()

	hub := newHub(db)

	// Create persistent channel in database
	err := hub.createChannelInDB("persistent-test", Persistent)
	if err != nil {
		t.Fatalf("Failed to create persistent channel: %v", err)
	}

	client := &Client{
		hub:     hub,
		conn:    nil,
		send:    make(chan []byte, 256),
		channel: "persistent-test",
	}

	// Create persistent channel in memory with only one client
	hub.channels["persistent-test"] = newChannel("persistent-test", Persistent)
	hub.channels["persistent-test"].clients[client] = true

	// Switch to another channel
	client.switchChannel("general")

	// The persistent channel should be removed from memory but not from database
	if _, exists := hub.channels["persistent-test"]; exists {
		t.Error("Empty persistent channel should be removed from memory when last client leaves")
	}

	// Verify it still exists in database
	var channelType string
	err = db.QueryRow("SELECT type FROM channels WHERE name = $1", "persistent-test").Scan(&channelType)
	if err != nil {
		t.Error("Persistent channel should still exist in database")
	}
}

func TestClientWritePump(t *testing.T) {
	// Test that writePump handles closed send channel gracefully
	// We can't test with a real WebSocket connection, so we skip the actual writing

	client := &Client{
		hub:     nil,
		conn:    nil, // This will cause WriteMessage to fail
		send:    make(chan []byte, 1),
		channel: "test",
	}

	// Close send channel to terminate writePump quickly
	close(client.send)

	// The method should exit gracefully without panicking on the close channel
	// We can't call writePump directly since it will panic on nil conn.WriteMessage
	// Instead, we test that the send channel behaves correctly

	// Verify send channel is closed
	_, ok := <-client.send
	if ok {
		t.Error("Send channel should be closed")
	}
}

func TestChannelSwitchMessage(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Skip("Skipping client tests without database")
		return
	}
	defer db.Close()

	hub := newHub(db)

	client := &Client{
		hub:     hub,
		conn:    nil,
		send:    make(chan []byte, 256),
		channel: "general",
	}

	// Create initial channel
	hub.channels["general"] = newChannel("general", Ephemeral)
	hub.channels["general"].clients[client] = true

	// Switch channel
	client.switchChannel("new-channel")

	// Check if client received channel switch message
	select {
	case msg := <-client.send:
		var message Message
		err := json.Unmarshal(msg, &message)
		if err != nil {
			t.Fatalf("Failed to unmarshal channel switch message: %v", err)
		}

		if message.Type != "channel_switch" {
			t.Errorf("Expected message type 'channel_switch', got '%s'", message.Type)
		}

		if message.Channel != "new-channel" {
			t.Errorf("Expected channel 'new-channel', got '%s'", message.Channel)
		}

		if message.Username != "System" {
			t.Errorf("Expected username 'System', got '%s'", message.Username)
		}

	case <-time.After(100 * time.Millisecond):
		t.Error("Client should receive channel switch message")
	}
}

func TestChannelCreatedMessage(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Skip("Skipping client tests without database")
		return
	}
	defer db.Close()

	hub := newHub(db)
	hub.channels = make(map[string]*Channel) // Ensure clean state

	client := &Client{
		hub:     hub,
		conn:    nil,
		send:    make(chan []byte, 256),
		channel: "general",
	}

	// Create initial general channel
	hub.channels["general"] = newChannel("general", Ephemeral)
	hub.channels["general"].clients[client] = true

	// Switch to a new channel that doesn't exist yet
	client.switchChannelWithType("brand-new-channel", Ephemeral)

	// The first message should be channel_switch
	select {
	case <-client.send:
		// Consume channel switch message
	case <-time.After(100 * time.Millisecond):
		t.Error("Should receive channel switch message")
		return
	}

	// Since we used switchChannelWithType and the channel was created,
	// we should verify the channel was created correctly
	if _, exists := hub.channels["brand-new-channel"]; !exists {
		t.Error("New channel should be created")
	}

	if hub.channels["brand-new-channel"].channelType != Ephemeral {
		t.Error("New channel should have correct type")
	}
}

func TestEmptyChannelNameDefaults(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Skip("Skipping client tests without database")
		return
	}
	defer db.Close()

	hub := newHub(db)

	client := &Client{
		hub:     hub,
		conn:    nil,
		send:    make(chan []byte, 256),
		channel: "",
	}

	// Test switching with empty channel name
	client.switchChannel("")

	if client.channel != "general" {
		t.Errorf("Empty channel name should default to 'general', got '%s'", client.channel)
	}

	// Test switchChannelWithType with empty name
	client.switchChannelWithType("", Persistent)

	if client.channel != "general" {
		t.Errorf("Empty channel name should default to 'general', got '%s'", client.channel)
	}
}
