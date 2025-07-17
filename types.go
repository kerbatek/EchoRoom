package main

import (
	"database/sql"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	channels   map[string]*Channel
	channelsMu sync.RWMutex
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	db         *sql.DB
	shutdown   chan bool
}

type ChannelType string

const (
	Ephemeral  ChannelType = "ephemeral"
	Persistent ChannelType = "persistent"
)

type Channel struct {
	name        string
	channelType ChannelType
	clients     map[*Client]bool
	clientsMu   sync.RWMutex
	broadcast   chan []byte
	shutdown    chan bool
}

type Client struct {
	hub       *Hub
	conn      *websocket.Conn
	send      chan []byte
	channel   string
	username  string
	hasJoined bool
}

type Message struct {
	ID        int       `json:"id,omitempty"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	Channel   string    `json:"channel"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

type ChannelCreateRequest struct {
	Name        string      `json:"name"`
	ChannelType ChannelType `json:"channel_type"`
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}
