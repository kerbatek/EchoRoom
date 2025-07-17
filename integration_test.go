package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestFullWorkflow(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Skip("Skipping integration test without database")
		return
	}
	defer db.Close()

	// Setup server
	hub := newHub(db)
	go hub.run()
	defer hub.stop() // Ensure cleanup

	// Reset ServeMux for clean test
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(hub, w, r)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	// Connect first client
	conn1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect first client: %v", err)
	}
	defer conn1.Close()

	// Connect second client
	conn2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect second client: %v", err)
	}
	defer conn2.Close()

	// Give time for connections to register
	time.Sleep(100 * time.Millisecond)

	// Drain initial active_channels messages
	conn1.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	conn1.ReadMessage() // Drain active_channels
	conn2.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	conn2.ReadMessage() // Drain active_channels

	// Reset deadlines
	conn1.SetReadDeadline(time.Time{})
	conn2.SetReadDeadline(time.Time{})

	// Test 1: Create a persistent channel
	createChannelMsg := ChannelCreateRequest{
		Name:        "test-persistent",
		ChannelType: Persistent,
	}
	createMsgBytes, _ := json.Marshal(map[string]interface{}{
		"type":         "create_channel",
		"name":         createChannelMsg.Name,
		"channel_type": createChannelMsg.ChannelType,
	})

	err = conn1.WriteMessage(websocket.TextMessage, createMsgBytes)
	if err != nil {
		t.Fatalf("Failed to send create channel message: %v", err)
	}

	// Read messages from client 1 (might get channel_created first, then channel_switch)
	found := false
	for i := 0; i < 3; i++ {
		_, msg, err := conn1.ReadMessage()
		if err != nil {
			t.Fatalf("Failed to read message: %v", err)
		}

		var tempMsg Message
		json.Unmarshal(msg, &tempMsg)
		if tempMsg.Type == "channel_switch" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected channel_switch message")
	}

	// Read channel created broadcast from client 2
	conn2.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	_, _, err = conn2.ReadMessage()
	if err != nil {
		// Might be active_channels message first, try again
		_, _, err = conn2.ReadMessage()
		if err != nil {
			t.Logf("Warning: Could not read message from conn2: %v", err)
		}
	}

	// Test 2: Send a message in the persistent channel
	chatMsg := Message{
		Username: "testuser1",
		Content:  "Hello from persistent channel!",
		Type:     "message",
		Channel:  "test-persistent",
	}

	chatMsgBytes, _ := json.Marshal(chatMsg)
	err = conn1.WriteMessage(websocket.TextMessage, chatMsgBytes)
	if err != nil {
		t.Fatalf("Failed to send chat message: %v", err)
	}

	// Test 3: Switch client 2 to the persistent channel
	switchChannelMsg := Message{
		Type:    "join_channel",
		Channel: "test-persistent",
	}

	switchMsgBytes, _ := json.Marshal(switchChannelMsg)
	err = conn2.WriteMessage(websocket.TextMessage, switchMsgBytes)
	if err != nil {
		t.Fatalf("Failed to send switch channel message: %v", err)
	}

	// Client 2 should receive:
	// 1. Channel switch confirmation
	// 2. Message history (including the message sent by client 1)

	messagesReceived := 0
	for messagesReceived < 3 { // switch msg + history msg + any other messages
		conn2.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		_, msg, err := conn2.ReadMessage()
		if err != nil {
			break
		}

		var receivedMsg Message
		if json.Unmarshal(msg, &receivedMsg) == nil {
			if receivedMsg.Type == "channel_switch" {
				if receivedMsg.Channel != "test-persistent" {
					t.Errorf("Expected switch to test-persistent, got %s", receivedMsg.Channel)
				}
				messagesReceived++
			} else if receivedMsg.Type == "message" {
				if receivedMsg.Content == "Hello from persistent channel!" {
					messagesReceived++
				}
			}
		}
		messagesReceived++
		if messagesReceived >= 5 { // Prevent infinite loop
			break
		}
	}

	// Test 4: Send a message from client 2, client 1 should receive it
	chatMsg2 := Message{
		Username: "testuser2",
		Content:  "Hello from client 2!",
		Type:     "message",
		Channel:  "test-persistent",
	}

	chatMsg2Bytes, _ := json.Marshal(chatMsg2)
	err = conn2.WriteMessage(websocket.TextMessage, chatMsg2Bytes)
	if err != nil {
		t.Fatalf("Failed to send second chat message: %v", err)
	}

	// Client 1 should receive the message (might get system messages first)
	var receivedChatMsg Message
	found = false
	for i := 0; i < 5; i++ {
		conn1.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		_, msg, err := conn1.ReadMessage()
		if err != nil {
			t.Fatalf("Client 1 should receive message from client 2: %v", err)
		}

		err = json.Unmarshal(msg, &receivedChatMsg)
		if err != nil {
			continue
		}

		if receivedChatMsg.Type == "message" && receivedChatMsg.Content == "Hello from client 2!" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected 'Hello from client 2!', got '%s'", receivedChatMsg.Content)
	}

	// Test 5: Verify message persistence
	// Disconnect both clients and reconnect one to check message history
	conn1.Close()
	conn2.Close()

	time.Sleep(100 * time.Millisecond)

	// Reconnect client
	conn3, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to reconnect client: %v", err)
	}
	defer conn3.Close()

	// Switch to persistent channel
	err = conn3.WriteMessage(websocket.TextMessage, switchMsgBytes)
	if err != nil {
		t.Fatalf("Failed to send switch message: %v", err)
	}

	// Should receive channel switch + message history
	historyReceived := false
	for i := 0; i < 10; i++ { // Try multiple reads
		conn3.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
		_, msg, err := conn3.ReadMessage()
		if err != nil {
			break
		}

		var historyMsg Message
		if json.Unmarshal(msg, &historyMsg) == nil {
			if historyMsg.Type == "message" &&
				(historyMsg.Content == "Hello from persistent channel!" ||
					historyMsg.Content == "Hello from client 2!") {
				historyReceived = true
				break
			}
		}
	}

	if !historyReceived {
		t.Error("Should receive message history when joining persistent channel")
	}
}

func TestEphemeralChannelCleanup(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Skip("Skipping ephemeral channel test without database")
		return
	}
	defer db.Close()

	hub := newHub(db)
	go hub.run()
	defer hub.stop() // Ensure cleanup

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(hub, w, r)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	// Connect client
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Create ephemeral channel
	createEphemeralMsg := map[string]interface{}{
		"type":         "create_channel",
		"name":         "ephemeral-test",
		"channel_type": "ephemeral",
	}

	msgBytes, _ := json.Marshal(createEphemeralMsg)
	err = conn.WriteMessage(websocket.TextMessage, msgBytes)
	if err != nil {
		t.Fatalf("Failed to create ephemeral channel: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Verify channel exists
	if _, exists := hub.channels["ephemeral-test"]; !exists {
		t.Error("Ephemeral channel should exist after creation")
	}

	// Disconnect client
	conn.Close()
	time.Sleep(200 * time.Millisecond)

	// Verify ephemeral channel was cleaned up
	if _, exists := hub.channels["ephemeral-test"]; exists {
		t.Error("Ephemeral channel should be cleaned up after last client disconnects")
	}

	// Verify it's not in database
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM channels WHERE name = 'ephemeral-test'").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query database: %v", err)
	}
	if count > 0 {
		t.Error("Ephemeral channel should not be stored in database")
	}
}

func TestConcurrentClients(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Skip("Skipping concurrent clients test without database")
		return
	}
	defer db.Close()

	hub := newHub(db)
	go hub.run()
	defer hub.stop() // Ensure cleanup

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(hub, w, r)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	// Connect multiple clients
	numClients := 5
	connections := make([]*websocket.Conn, numClients)

	for i := 0; i < numClients; i++ {
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			t.Fatalf("Failed to connect client %d: %v", i, err)
		}
		connections[i] = conn
		defer conn.Close()
	}

	time.Sleep(100 * time.Millisecond)

	// Give time for all initial messages to settle
	time.Sleep(200 * time.Millisecond)

	// Reset deadlines
	for i := 0; i < numClients; i++ {
		connections[i].SetReadDeadline(time.Time{})
	}

	// All clients should be in general channel
	generalChannel, exists := hub.channels["general"]
	if !exists {
		t.Fatal("General channel should exist")
	}

	if len(generalChannel.clients) != numClients {
		t.Errorf("Expected %d clients in general channel, got %d", numClients, len(generalChannel.clients))
	}

	// Test broadcasting to all clients
	testMsg := Message{
		Username: "broadcaster",
		Content:  "Hello everyone!",
		Type:     "message",
		Channel:  "general",
	}

	msgBytes, _ := json.Marshal(testMsg)
	err := connections[0].WriteMessage(websocket.TextMessage, msgBytes)
	if err != nil {
		t.Fatalf("Failed to send broadcast message: %v", err)
	}

	// All other clients should receive the message
	for i := 1; i < numClients; i++ {
		found := false
		connections[i].SetReadDeadline(time.Now().Add(1000 * time.Millisecond))

		// Keep reading until we find the broadcast message we're looking for
		for !found {
			_, msg, err := connections[i].ReadMessage()
			if err != nil {
				t.Errorf("Client %d did not receive broadcast message: %v", i, err)
				break
			}

			var receivedMsg Message
			if json.Unmarshal(msg, &receivedMsg) == nil {
				if receivedMsg.Type == "message" && receivedMsg.Content == "Hello everyone!" {
					found = true
				}
			}
		}

		if !found {
			t.Errorf("Client %d did not receive the expected broadcast message", i)
		}
	}
}
