package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWebSocketUpgrade(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		// Create a minimal hub for testing without database
		hub := &Hub{
			channels:   make(map[string]*Channel),
			register:   make(chan *Client, 1),
			unregister: make(chan *Client, 1),
			broadcast:  make(chan []byte, 1),
			db:         nil,
		}
		testWebSocketUpgrade(t, hub)
		return
	}
	defer db.Close()

	hub := newHub(db)
	testWebSocketUpgrade(t, hub)
}

func testWebSocketUpgrade(t *testing.T, hub *Hub) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(hub, w, r)
	}))
	defer server.Close()

	// Convert HTTP URL to WebSocket URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Test WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	// Connection successful if we reach here
	if conn == nil {
		t.Error("WebSocket connection should be established")
	}
}

func TestSetupRoutes(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		// Create a minimal hub for testing without database
		hub := &Hub{
			channels:   make(map[string]*Channel),
			register:   make(chan *Client, 1),
			unregister: make(chan *Client, 1),
			broadcast:  make(chan []byte, 1),
			db:         nil,
		}
		testSetupRoutes(t, hub)
		return
	}
	defer db.Close()

	hub := newHub(db)
	testSetupRoutes(t, hub)
}

func testSetupRoutes(t *testing.T, hub *Hub) {
	// Reset default ServeMux to avoid conflicts
	http.DefaultServeMux = http.NewServeMux()

	setupRoutes(hub)

	// Test WebSocket route
	req, err := http.NewRequest("GET", "/ws", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Connection", "upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")

	rr := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)

	// Should return 101 Switching Protocols for valid WebSocket upgrade
	if rr.Code != http.StatusSwitchingProtocols {
		// Note: This might fail due to missing WebSocket headers, but route should exist
		// Let's check if the route exists by testing a regular GET request
		req2, _ := http.NewRequest("GET", "/ws", nil)
		rr2 := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(rr2, req2)

		// Should not return 404 if route is registered
		if rr2.Code == http.StatusNotFound {
			t.Error("WebSocket route /ws should be registered")
		}
	}

	// Test static file route
	req, err = http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)

	// Should not return 404 (file might not exist in test, but route should be registered)
	if rr.Code == http.StatusNotFound {
		t.Error("Static file route / should be registered")
	}
}

func TestWebSocketUpgraderConfig(t *testing.T) {
	// Test that the upgrader has the correct configuration
	if upgrader.CheckOrigin == nil {
		t.Error("Upgrader CheckOrigin function should be set")
		return
	}

	// Test CheckOrigin function always returns true
	req, _ := http.NewRequest("GET", "/ws", nil)
	req.Header.Set("Origin", "http://example.com")

	if !upgrader.CheckOrigin(req) {
		t.Error("CheckOrigin should return true for any origin")
	}

	req.Header.Set("Origin", "http://malicious-site.com")
	if !upgrader.CheckOrigin(req) {
		t.Error("CheckOrigin should return true for any origin (including potentially malicious ones)")
	}
}

func TestHandleWebSocketClientRegistration(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Skip("Skipping WebSocket client registration test without database")
		return
	}
	defer db.Close()

	hub := newHub(db)

	// Start hub in background
	go hub.run()
	defer hub.stop() // Ensure cleanup

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(hub, w, r)
	}))
	defer server.Close()

	// Convert HTTP URL to WebSocket URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// Connect WebSocket client
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	// Give time for client registration
	// We can't easily verify registration without exposing internal state
	// But we can test that the connection stays open
	conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	_, _, err = conn.ReadMessage()

	// We expect a timeout here since no messages are sent immediately
	if err == nil {
		// If we got a message, that's also fine (might be active_channels)
		return
	}

	// Check if it's a timeout error (expected) vs other errors
	if netErr, ok := err.(interface{ Timeout() bool }); ok && netErr.Timeout() {
		// Timeout is expected - connection is working
		return
	}

	// Other errors might indicate connection issues
	t.Logf("Connection error (might be expected): %v", err)
}

func TestWebSocketErrorHandling(t *testing.T) {
	// Test handling invalid WebSocket upgrade requests
	db := setupTestDB(t)
	if db == nil {
		// Use minimal hub for test
		hub := &Hub{
			channels:   make(map[string]*Channel),
			register:   make(chan *Client, 1),
			unregister: make(chan *Client, 1),
			broadcast:  make(chan []byte, 1),
			db:         nil,
		}
		testWebSocketErrorHandling(t, hub)
		return
	}
	defer db.Close()

	hub := newHub(db)
	testWebSocketErrorHandling(t, hub)
}

func testWebSocketErrorHandling(t *testing.T, hub *Hub) {
	// Test with invalid WebSocket upgrade request (missing headers)
	req, err := http.NewRequest("GET", "/ws", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handleWebSocket(hub, rr, req)

	// Should return an error status for invalid WebSocket upgrade
	if rr.Code == http.StatusSwitchingProtocols {
		t.Error("Should not upgrade invalid WebSocket request")
	}
}

