package main

import (
	"testing"
	"time"
)

func TestNewChannel(t *testing.T) {
	// Test creating ephemeral channel
	ephemeralChannel := newChannel("test-ephemeral", Ephemeral)

	if ephemeralChannel.name != "test-ephemeral" {
		t.Errorf("Expected channel name 'test-ephemeral', got '%s'", ephemeralChannel.name)
	}

	if ephemeralChannel.channelType != Ephemeral {
		t.Errorf("Expected channel type Ephemeral, got %s", ephemeralChannel.channelType)
	}

	if ephemeralChannel.clients == nil {
		t.Error("Channel clients map should be initialized")
	}

	if ephemeralChannel.broadcast == nil {
		t.Error("Channel broadcast channel should be initialized")
	}

	// Test creating persistent channel
	persistentChannel := newChannel("test-persistent", Persistent)

	if persistentChannel.name != "test-persistent" {
		t.Errorf("Expected channel name 'test-persistent', got '%s'", persistentChannel.name)
	}

	if persistentChannel.channelType != Persistent {
		t.Errorf("Expected channel type Persistent, got %s", persistentChannel.channelType)
	}
}

func TestChannelBroadcast(t *testing.T) {
	channel := newChannel("test-channel", Ephemeral)

	// Create mock clients
	client1 := &Client{
		hub:     nil,
		conn:    nil,
		send:    make(chan []byte, 256),
		channel: "test-channel",
	}

	client2 := &Client{
		hub:     nil,
		conn:    nil,
		send:    make(chan []byte, 256),
		channel: "test-channel",
	}

	// Add clients to channel
	channel.clients[client1] = true
	channel.clients[client2] = true

	// Start channel
	shutdown := make(chan bool)
	defer func() { shutdown <- true }() // Ensure cleanup
	go channel.run(shutdown)

	// Test broadcasting a message
	testMessage := []byte("test broadcast message")

	go func() {
		channel.broadcast <- testMessage
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

func TestChannelClientManagement(t *testing.T) {
	channel := newChannel("test-channel", Ephemeral)

	// Test initial state
	if len(channel.clients) != 0 {
		t.Error("New channel should have no clients")
	}

	// Add a client
	client := &Client{
		hub:     nil,
		conn:    nil,
		send:    make(chan []byte, 256),
		channel: "test-channel",
	}

	channel.clients[client] = true

	if len(channel.clients) != 1 {
		t.Error("Channel should have 1 client after adding")
	}

	if !channel.clients[client] {
		t.Error("Client should be present in channel")
	}

	// Remove client
	delete(channel.clients, client)

	if len(channel.clients) != 0 {
		t.Error("Channel should have no clients after removal")
	}
}

func TestChannelBroadcastWithBlockedClient(t *testing.T) {
	channel := newChannel("test-channel", Ephemeral)

	// Create a client with no buffer
	client := &Client{
		hub:     nil,
		conn:    nil,
		send:    make(chan []byte), // No buffer - will block immediately
		channel: "test-channel",
	}

	// Add client to channel
	channel.clients[client] = true

	// Start channel
	shutdown := make(chan bool)
	defer func() { shutdown <- true }() // Ensure cleanup
	go channel.run(shutdown)

	// Give channel time to start
	time.Sleep(10 * time.Millisecond)

	// Test broadcasting a message to a client with no receiver
	testMessage := []byte("test message")

	// This should trigger the default case in channel.run() and remove the client
	done := make(chan bool)
	go func() {
		channel.broadcast <- testMessage
		done <- true
	}()

	// Give time for broadcast to process
	select {
	case <-done:
		// Broadcast completed
	case <-time.After(100 * time.Millisecond):
		// Broadcast should complete even with blocked client
	}

	// Give additional time for client cleanup
	time.Sleep(50 * time.Millisecond)

	// Client should be removed from channel due to blocked send
	if len(channel.clients) != 0 {
		t.Error("Client with blocked send channel should be removed from channel")
	}
}

func TestChannelTypes(t *testing.T) {
	// Test that ChannelType constants are properly defined
	if Ephemeral != "ephemeral" {
		t.Errorf("Expected Ephemeral to be 'ephemeral', got '%s'", Ephemeral)
	}

	if Persistent != "persistent" {
		t.Errorf("Expected Persistent to be 'persistent', got '%s'", Persistent)
	}

	// Test channel type assignment
	ephemeralChannel := newChannel("test", Ephemeral)
	if ephemeralChannel.channelType != Ephemeral {
		t.Error("Ephemeral channel should have Ephemeral type")
	}

	persistentChannel := newChannel("test", Persistent)
	if persistentChannel.channelType != Persistent {
		t.Error("Persistent channel should have Persistent type")
	}
}

func TestChannelMultipleBroadcasts(t *testing.T) {
	channel := newChannel("test-channel", Ephemeral)

	// Create a client
	client := &Client{
		hub:     nil,
		conn:    nil,
		send:    make(chan []byte, 256),
		channel: "test-channel",
	}

	channel.clients[client] = true

	// Start channel
	shutdown := make(chan bool)
	defer func() { shutdown <- true }() // Ensure cleanup
	go channel.run(shutdown)

	// Send multiple messages
	messages := []string{"message1", "message2", "message3"}

	for _, msg := range messages {
		go func(m string) {
			channel.broadcast <- []byte(m)
		}(msg)
	}

	// Receive all messages
	receivedMessages := make([]string, 0, len(messages))
	for i := 0; i < len(messages); i++ {
		select {
		case msg := <-client.send:
			receivedMessages = append(receivedMessages, string(msg))
		case <-time.After(100 * time.Millisecond):
			t.Errorf("Failed to receive message %d", i+1)
		}
	}

	// Verify all messages were received (order may vary due to goroutines)
	if len(receivedMessages) != len(messages) {
		t.Errorf("Expected %d messages, received %d", len(messages), len(receivedMessages))
	}

	// Check that all original messages are present
	for _, originalMsg := range messages {
		found := false
		for _, receivedMsg := range receivedMessages {
			if receivedMsg == originalMsg {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Message '%s' was not received", originalMsg)
		}
	}
}
