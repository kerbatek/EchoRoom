package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func newHub(db *sql.DB) *Hub {
	return &Hub{
		channels:   make(map[string]*Channel),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		db:         db,
		shutdown:   make(chan bool),
	}
}

func (h *Hub) sendActiveChannels(client *Client) {
	// Get all channels from database (persistent channels)
	rows, err := h.db.Query("SELECT name, type FROM channels ORDER BY name")
	if err != nil {
		log.Printf("Error querying channels: %v", err)
		return
	}
	defer rows.Close()

	type ChannelInfo struct {
		Name string      `json:"name"`
		Type ChannelType `json:"type"`
	}

	var channelInfos []ChannelInfo
	channelMap := make(map[string]ChannelType)

	// Add persistent channels from database
	for rows.Next() {
		var name, channelType string
		if err := rows.Scan(&name, &channelType); err != nil {
			continue
		}
		channelInfos = append(channelInfos, ChannelInfo{Name: name, Type: ChannelType(channelType)})
		channelMap[name] = ChannelType(channelType)
	}

	// Add currently active ephemeral channels not in database
	h.channelsMu.RLock()
	for name, channel := range h.channels {
		if _, exists := channelMap[name]; !exists && channel.channelType == Ephemeral {
			channelInfos = append(channelInfos, ChannelInfo{Name: name, Type: Ephemeral})
		}
	}
	h.channelsMu.RUnlock()

	activeChannelsMsg := struct {
		Type     string        `json:"type"`
		Channels []ChannelInfo `json:"channels"`
	}{
		Type:     "active_channels",
		Channels: channelInfos,
	}

	if msgBytes, err := json.Marshal(activeChannelsMsg); err == nil {
		select {
		case client.send <- msgBytes:
		default:
			close(client.send)
		}
	}
}

func (h *Hub) run() {
	for {
		select {
		case <-h.shutdown:
			return
		case client := <-h.register:
			channelName := client.channel
			if channelName == "" {
				channelName = "general"
			}

			h.channelsMu.Lock()
			if _, ok := h.channels[channelName]; !ok {
				// Get channel type from database, default to ephemeral if not found
				channelType, err := h.getChannelType(channelName)
				if err != nil {
					channelType = Ephemeral
				}

				h.channels[channelName] = newChannel(channelName, channelType)
				go h.channels[channelName].run(h.shutdown)
				// Give the channel goroutine a moment to start
				time.Sleep(1 * time.Millisecond)
			}
			channel := h.channels[channelName]
			h.channelsMu.Unlock()

			channel.clientsMu.Lock()
			channel.clients[client] = true
			clientCount := len(channel.clients)
			channel.clientsMu.Unlock()
			log.Printf("Client connected to channel '%s'. Total clients in channel: %d", channelName, clientCount)

			// Send message history for persistent channels
			if channel.channelType == Persistent {
				history, err := h.getChannelHistory(channelName, 50) // Last 50 messages
				if err == nil {
					for _, msg := range history {
						if msgBytes, err := json.Marshal(msg); err == nil {
							select {
							case client.send <- msgBytes:
							default:
								close(client.send)
								return
							}
						}
					}
				}
			}

			// Send active channels list to the newly connected client
			h.sendActiveChannels(client)

		case client := <-h.unregister:
			channelName := client.channel
			if channelName == "" {
				channelName = "general"
			}

			h.channelsMu.RLock()
			channel, ok := h.channels[channelName]
			h.channelsMu.RUnlock()

			if ok {
				channel.clientsMu.Lock()
				if _, ok := channel.clients[client]; ok {
					// Send leave message for ephemeral channels only if there will be other clients remaining
					if channel.channelType == Ephemeral && len(channel.clients) > 1 && client.username != "" {
						leaveMsg := Message{
							Username:  "System",
							Content:   fmt.Sprintf("%s left the channel", client.username),
							Type:      "system_message",
							Channel:   channelName,
							Timestamp: time.Now().UTC(),
						}
						if leaveMsgBytes, err := json.Marshal(leaveMsg); err == nil {
							log.Printf("Sending leave message for %s to channel '%s'", client.username, channelName)
							channel.broadcast <- leaveMsgBytes
						}
					}

					delete(channel.clients, client)
					clientCount := len(channel.clients)
					channel.clientsMu.Unlock()
					log.Printf("Client disconnected from channel '%s'. Total clients in channel: %d", channelName, clientCount)

					if clientCount == 0 && channelName != "general" {
						// Only remove ephemeral channels when empty
						if channel.channelType == Ephemeral {
							h.channelsMu.Lock()
							delete(h.channels, channelName)
							h.channelsMu.Unlock()
							log.Printf("Ephemeral channel '%s' removed (no clients)", channelName)

							// Broadcast channel deletion to all clients BEFORE closing the send channel
							channelDeletedMsg := Message{
								Username: "System",
								Content:  channelName,
								Type:     "channel_deleted",
								Channel:  channelName,
							}

							if msgBytes, err := json.Marshal(channelDeletedMsg); err == nil {
								// Send to all remaining clients in all channels
								h.channelsMu.RLock()
								for _, ch := range h.channels {
									ch.clientsMu.RLock()
									for c := range ch.clients {
										if c != client { // Don't send to the disconnecting client
											select {
											case c.send <- msgBytes:
											default:
												// Client's send channel is full, skip
											}
										}
									}
									ch.clientsMu.RUnlock()
								}
								h.channelsMu.RUnlock()
							}
						} else {
							// For persistent channels, keep the channel but remove it from memory
							h.channelsMu.Lock()
							delete(h.channels, channelName)
							h.channelsMu.Unlock()
							log.Printf("Persistent channel '%s' removed from memory (no clients, but preserved in database)", channelName)
						}
					}

					close(client.send)
				}
			}
		case message := <-h.broadcast:
			h.channelsMu.RLock()
			for _, channel := range h.channels {
				channel.clientsMu.Lock()
				for client := range channel.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(channel.clients, client)
					}
				}
				channel.clientsMu.Unlock()
			}
			h.channelsMu.RUnlock()
		}
	}
}

func (h *Hub) stop() {
	// Stop all channels first
	h.channelsMu.RLock()
	for _, channel := range h.channels {
		select {
		case channel.shutdown <- true:
		default:
		}
	}
	h.channelsMu.RUnlock()

	// Stop the hub
	select {
	case h.shutdown <- true:
	default:
	}
}
