package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	
	"github.com/gorilla/websocket"
)

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// First, check the message type to determine how to unmarshal
		var msgType struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(messageBytes, &msgType); err != nil {
			log.Printf("Error unmarshaling message type: %v", err)
			continue
		}

		if msgType.Type == "user_connected" {
			var message Message
			if err := json.Unmarshal(messageBytes, &message); err != nil {
				log.Printf("Error unmarshaling user_connected message: %v", err)
				continue
			}
			
			// Set username and send join message immediately
			if message.Username != "" && c.username != message.Username {
				c.username = message.Username
				c.hasJoined = true
				
				channelName := c.channel
				if channelName == "" {
					channelName = "general"
				}
				
				// Send join message for ephemeral channels if there are other clients
				if channel, ok := c.hub.channels[channelName]; ok && channel.channelType == Ephemeral {
					if len(channel.clients) > 1 {
						joinMsg := Message{
							Username:  "System",
							Content:   fmt.Sprintf("%s joined the channel", c.username),
							Type:      "system_message",
							Channel:   channelName,
							Timestamp: time.Now().UTC(),
						}
						if joinMsgBytes, err := json.Marshal(joinMsg); err == nil {
							log.Printf("Sending immediate join message for %s to channel '%s'", c.username, channelName)
							channel.broadcast <- joinMsgBytes
						}
					}
				}
			}
			continue
		}

		if msgType.Type == "join_channel" {
			var message Message
			if err := json.Unmarshal(messageBytes, &message); err != nil {
				log.Printf("Error unmarshaling join_channel message: %v", err)
				continue
			}
			c.switchChannel(message.Channel)
			continue
		}

		if msgType.Type == "create_channel" {
			var createReq ChannelCreateRequest
			if err := json.Unmarshal(messageBytes, &createReq); err != nil {
				log.Printf("Error unmarshaling channel create request: %v", err)
				continue
			}

			log.Printf("Received create_channel request: name='%s', channel_type='%s'", createReq.Name, createReq.ChannelType)

			// Create channel in database (only for persistent channels)
			if err := c.hub.createChannelInDB(createReq.Name, createReq.ChannelType); err != nil {
				log.Printf("Error creating channel in database: %v", err)
				continue
			}
			if createReq.ChannelType == Persistent {
				log.Printf("Persistent channel created in database: name='%s'", createReq.Name)
			} else {
				log.Printf("Ephemeral channel created: name='%s' (not stored in database)", createReq.Name)
			}

			// Switch to the new channel
			c.switchChannelWithType(createReq.Name, createReq.ChannelType)
			log.Printf("Switched client to new channel: %s", createReq.Name)
			continue
		}

		// Handle regular messages
		var message Message
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		channelName := c.channel
		if channelName == "" {
			channelName = "general"
		}

		// Update client username from message (username should already be set from user_connected)
		if message.Username != "" && c.username != message.Username {
			c.username = message.Username
		}

		// Only process regular messages for channel broadcasting
		if message.Type == "message" {
			// Set timestamp for all messages
			message.Timestamp = time.Now().UTC()
			
			// Only save to database if channel is persistent
			if channel, ok := c.hub.channels[channelName]; ok && channel.channelType == Persistent {
				if err := c.hub.saveMessage(message); err != nil {
					log.Printf("Error saving message: %v", err)
				}
			}

			// Broadcast to channel with updated timestamp
			if channel, ok := c.hub.channels[channelName]; ok {
				// Re-marshal the message with the timestamp included
				if updatedBytes, err := json.Marshal(message); err == nil {
					channel.broadcast <- updatedBytes
				} else {
					// Fallback to original message if marshaling fails
					channel.broadcast <- messageBytes
				}
			}
		}
	}
}

func (c *Client) switchChannelWithType(newChannelName string, channelType ChannelType) {
	if newChannelName == "" {
		newChannelName = "general"
		channelType = Ephemeral
	}

	if c.channel == newChannelName {
		return
	}

	oldChannel := c.channel
	if oldChannel == "" {
		oldChannel = "general"
	}

	if channel, ok := c.hub.channels[oldChannel]; ok {
		delete(channel.clients, c)
		if len(channel.clients) == 0 && oldChannel != "general" {
			// Only delete ephemeral channels when empty
			if channel.channelType == Ephemeral {
				delete(c.hub.channels, oldChannel)

				// Broadcast channel deletion to all clients
				channelDeletedMsg := Message{
					Username: "System",
					Content:  oldChannel,
					Type:     "channel_deleted",
					Channel:  oldChannel,
				}

				if msgBytes, err := json.Marshal(channelDeletedMsg); err == nil {
					select {
					case c.hub.broadcast <- msgBytes:
					default:
						// Hub broadcast channel is full, skip
					}
				}
			} else {
				// For persistent channels, just remove from memory
				delete(c.hub.channels, oldChannel)
			}
		}
	}

	c.channel = newChannelName

	channelCreated := false
	if _, ok := c.hub.channels[newChannelName]; !ok {
		c.hub.channels[newChannelName] = newChannel(newChannelName, channelType)
		go c.hub.channels[newChannelName].run(c.hub.shutdown)
		channelCreated = true
	}

	c.hub.channels[newChannelName].clients[c] = true

	// Send join message for ephemeral channels if there are other clients and we have a username
	if newChannel := c.hub.channels[newChannelName]; newChannel.channelType == Ephemeral && len(newChannel.clients) > 1 && c.username != "" {
		joinMsg := Message{
			Username:  "System",
			Content:   fmt.Sprintf("%s joined the channel", c.username),
			Type:      "system_message",
			Channel:   newChannelName,
			Timestamp: time.Now().UTC(),
		}
		if joinMsgBytes, err := json.Marshal(joinMsg); err == nil {
			log.Printf("Sending join message for %s switching to channel '%s'", c.username, newChannelName)
			newChannel.broadcast <- joinMsgBytes
		}
	}

	if channelCreated {
		channelCreatedMsg := struct {
			Type        string      `json:"type"`
			Name        string      `json:"name"`
			ChannelType ChannelType `json:"channel_type"`
		}{
			Type:        "channel_created",
			Name:        newChannelName,
			ChannelType: channelType,
		}

		if msgBytes, err := json.Marshal(channelCreatedMsg); err == nil {
			select {
		case c.hub.broadcast <- msgBytes:
		default:
			// Hub broadcast channel is full, skip
		}
		}
	}

	channelSwitchMsg := Message{
		Username: "System",
		Content:  "Switched to channel: " + newChannelName,
		Type:     "channel_switch",
		Channel:  newChannelName,
	}

	if msgBytes, err := json.Marshal(channelSwitchMsg); err == nil {
		c.send <- msgBytes
	}

	// Send message history for persistent channels AFTER channel switch message
	if c.hub.channels[newChannelName].channelType == Persistent {
		history, err := c.hub.getChannelHistory(newChannelName, 50)
		if err == nil {
			log.Printf("Loading %d messages from history for channel '%s'", len(history), newChannelName)
			for _, msg := range history {
				if msgBytes, err := json.Marshal(msg); err == nil {
					select {
					case c.send <- msgBytes:
					default:
						close(c.send)
						return
					}
				}
			}
		} else {
			log.Printf("Error loading message history for channel '%s': %v", newChannelName, err)
		}
	}

	log.Printf("Client switched from '%s' to '%s'", oldChannel, newChannelName)
}

func (c *Client) switchChannel(newChannelName string) {
	if newChannelName == "" {
		newChannelName = "general"
	}

	if c.channel == newChannelName {
		return
	}

	oldChannel := c.channel
	if oldChannel == "" {
		oldChannel = "general"
	}

	if channel, ok := c.hub.channels[oldChannel]; ok {
		// Send leave message for ephemeral channels if there are other clients
		if channel.channelType == Ephemeral && len(channel.clients) > 1 && c.username != "" {
			leaveMsg := Message{
				Username:  "System",
				Content:   fmt.Sprintf("%s left the channel", c.username),
				Type:      "system_message",
				Channel:   oldChannel,
				Timestamp: time.Now().UTC(),
			}
			if leaveMsgBytes, err := json.Marshal(leaveMsg); err == nil {
				log.Printf("Sending leave message for %s switching from channel '%s'", c.username, oldChannel)
				channel.broadcast <- leaveMsgBytes
			}
		}
		
		delete(channel.clients, c)
		if len(channel.clients) == 0 && oldChannel != "general" {
			// Only delete ephemeral channels when empty
			if channel.channelType == Ephemeral {
				delete(c.hub.channels, oldChannel)

				// Broadcast channel deletion to all clients
				channelDeletedMsg := Message{
					Username: "System",
					Content:  oldChannel,
					Type:     "channel_deleted",
					Channel:  oldChannel,
				}

				if msgBytes, err := json.Marshal(channelDeletedMsg); err == nil {
					select {
					case c.hub.broadcast <- msgBytes:
					default:
						// Hub broadcast channel is full, skip
					}
				}
			} else {
				// For persistent channels, just remove from memory
				delete(c.hub.channels, oldChannel)
			}
		}
	}

	c.channel = newChannelName

	channelCreated := false
	if _, ok := c.hub.channels[newChannelName]; !ok {
		// Get channel type from database, default to ephemeral if not found
		channelType, err := c.hub.getChannelType(newChannelName)
		if err != nil {
			channelType = Ephemeral
		}

		c.hub.channels[newChannelName] = newChannel(newChannelName, channelType)
		go c.hub.channels[newChannelName].run(c.hub.shutdown)
		channelCreated = true
	}

	c.hub.channels[newChannelName].clients[c] = true

	// Send join message for ephemeral channels if there are other clients and we have a username
	if newChannel := c.hub.channels[newChannelName]; newChannel.channelType == Ephemeral && len(newChannel.clients) > 1 && c.username != "" {
		joinMsg := Message{
			Username:  "System",
			Content:   fmt.Sprintf("%s joined the channel", c.username),
			Type:      "system_message",
			Channel:   newChannelName,
			Timestamp: time.Now().UTC(),
		}
		if joinMsgBytes, err := json.Marshal(joinMsg); err == nil {
			log.Printf("Sending join message for %s switching to channel '%s'", c.username, newChannelName)
			newChannel.broadcast <- joinMsgBytes
		}
	}

	if channelCreated {
		// Get the channel type to send in the message
		channelType := c.hub.channels[newChannelName].channelType

		channelCreatedMsg := struct {
			Type        string      `json:"type"`
			Name        string      `json:"name"`
			ChannelType ChannelType `json:"channel_type"`
		}{
			Type:        "channel_created",
			Name:        newChannelName,
			ChannelType: channelType,
		}

		if msgBytes, err := json.Marshal(channelCreatedMsg); err == nil {
			select {
		case c.hub.broadcast <- msgBytes:
		default:
			// Hub broadcast channel is full, skip
		}
		}
	}

	channelSwitchMsg := Message{
		Username: "System",
		Content:  "Switched to channel: " + newChannelName,
		Type:     "channel_switch",
		Channel:  newChannelName,
	}

	if msgBytes, err := json.Marshal(channelSwitchMsg); err == nil {
		c.send <- msgBytes
	}

	// Send message history for persistent channels AFTER channel switch message
	if c.hub.channels[newChannelName].channelType == Persistent {
		history, err := c.hub.getChannelHistory(newChannelName, 50)
		if err == nil {
			log.Printf("Loading %d messages from history for channel '%s'", len(history), newChannelName)
			for _, msg := range history {
				if msgBytes, err := json.Marshal(msg); err == nil {
					select {
					case c.send <- msgBytes:
					default:
						close(c.send)
						return
					}
				}
			}
		} else {
			log.Printf("Error loading message history for channel '%s': %v", newChannelName, err)
		}
	}

	log.Printf("Client switched from '%s' to '%s'", oldChannel, newChannelName)
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for message := range c.send {
		c.conn.WriteMessage(websocket.TextMessage, message)
	}
	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
}