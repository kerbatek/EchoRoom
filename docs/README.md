# 🌍 EchoRoom - Real-time Chat Application

A modern WebSocket-based chat application with PostgreSQL persistence, featuring premium branding, desktop notifications, and intelligent user experience enhancements. Supports both ephemeral and persistent channels with comprehensive notification system.

## Quick Start

### 1. Database Setup
Ensure you have PostgreSQL running and create a database:
```sql
CREATE DATABASE chat_app;
```

### 2. Environment Configuration
Copy the example environment file and configure:
```bash
cp configs/.env.example .env
# Edit .env with your PostgreSQL connection details
```

Or set environment variables directly:
```bash
# Windows
set DB_HOST=localhost
set DB_USER=postgres
set DB_PASSWORD=your_password
set DB_NAME=chat_app

# Linux/Mac
export DB_HOST=localhost
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=chat_app
```

### 3. Run the Application
```bash
go run .
```

Open your browser to `http://localhost:8080`

## Testing

### Prerequisites
- PostgreSQL database running
- Test database configured (recommended: `chat_app_test`)

### Setup Test Environment
```bash
# Option 1: Use separate test database (recommended)
export TEST_DB_NAME=chat_app_test
createdb chat_app_test

# Option 2: Use same database as development
# Tests will automatically skip if PostgreSQL is not available
```

### Run Tests
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

**Note**: Database-dependent tests will be automatically skipped if PostgreSQL is not reachable. Core functionality tests will still run.

## ✨ Key Features

### 🚀 Core Chat Features
- **Ephemeral Channels** ⚡: Temporary channels that disappear when empty
- **Persistent Channels** 💾: Permanent channels with message history
- **Real-time Sync**: All clients see channel changes instantly
- **Message History**: Last 50 messages loaded for persistent channels
- **Smart Usernames**: 150+ funny, automatically assigned usernames

### 🎨 Premium User Experience
- **Modern Branding**: EchoRoom branding with gradient themes and premium typography
- **Dark/Light Theme**: Seamless theme toggle with localStorage persistence
- **Responsive Design**: Optimized for desktop and mobile devices
- **Loading Indicators**: Visual feedback for all network operations
- **Custom Favicon**: Branded chat bubble icon

### ⏰ Smart Timestamps
- **Relative Time**: User-friendly formats ("5 min ago", "just now")
- **Auto-Updates**: Timestamps refresh every 5 seconds
- **Tooltip Details**: Hover for full date/time information

### 🔔 Comprehensive Notifications
- **Desktop Notifications**: Browser alerts for new messages when tab inactive
- **Audio Alerts**: Pleasant beep sounds using Web Audio API
- **Title Blinking**: Slow blinking page title with unread message count
- **Permission Management**: Visual indicator for notification status
- **Smart Triggering**: Only notifies when page is not visible

## 🏠 Architecture

### 🔧 Backend (Modular Go)
- `main.go` - Application entry point and routing
- `types.go` - Type definitions (Hub, Channel, Client, Message)
- `hub.go` - Central hub managing clients and channels
- `channel.go` - Channel broadcast and lifecycle management
- `client.go` - WebSocket client connection handling
- `websocket.go` - WebSocket endpoint and upgrade logic
- `database.go` - PostgreSQL operations and schema management

### 🎨 Frontend (Modern Single Page App)
- `assets/index.html` - Premium EchoRoom branded interface
- `assets/styles.css` - Gradient themes with CSS custom properties
- `assets/script.js` - Advanced WebSocket client with notifications
- **Advanced Features**:
  - WebSocket connection with automatic reconnection and loading states
  - Real-time channel management and message display
  - Desktop notification system with permission management
  - Audio alerts using Web Audio API
  - Responsive design with sidebar and chat area
  - Smart timestamp formatting with auto-updates
  - Theme persistence and smooth transitions

### 🗺️ Database Schema (PostgreSQL)
**Channels Table:**
- `name` (VARCHAR(100) PRIMARY KEY) - Channel identifier
- `type` (VARCHAR(20) NOT NULL) - "ephemeral" or "persistent"
- `created_at` (TIMESTAMP DEFAULT CURRENT_TIMESTAMP) - Creation time

**Messages Table:**
- `id` (SERIAL PRIMARY KEY) - Auto-increment message ID
- `channel_name` (VARCHAR(100)) - References channels.name
- `username` (VARCHAR(100)) - Message author
- `content` (TEXT) - Message content
- `timestamp` (TIMESTAMP DEFAULT CURRENT_TIMESTAMP) - Message time

**Auto-Schema Creation**: Tables are automatically created on startup if they don't exist.

## 🚀 Commands

### Development
```bash
go run .                 # Start EchoRoom development server
go build -o echoroom .   # Build EchoRoom executable
go mod tidy              # Clean up dependencies
```

### Testing
```bash
go test ./...               # Run all tests
go test -v ./...           # Verbose test output
go test -cover ./...       # Test coverage report
```

### Production Deployment
```bash
# Build optimized binary
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o echoroom .

# Run with production settings
export DB_SSLMODE=require
./echoroom
```

## ⚙️ Environment Variables

### Database Configuration
| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | Database username |
| `DB_PASSWORD` | `password` | Database password |
| `DB_NAME` | `chat_app` | Database name |
| `DB_SSLMODE` | `disable` | SSL mode (disable/require) |

### Test Environment Variables
| Variable | Default | Description |
|----------|---------|-------------|
| `TEST_DB_HOST` | Same as `DB_HOST` | Test PostgreSQL host |
| `TEST_DB_PORT` | Same as `DB_PORT` | Test PostgreSQL port |
| `TEST_DB_USER` | Same as `DB_USER` | Test database username |
| `TEST_DB_PASSWORD` | Same as `DB_PASSWORD` | Test database password |
| `TEST_DB_NAME` | Same as `DB_NAME` | Test database name |
| `TEST_DB_SSLMODE` | Same as `DB_SSLMODE` | Test SSL mode |

## 🌐 Browser Support

**Required Features:**
- WebSocket support
- Notification API (optional, for desktop notifications)
- Web Audio API (optional, for sound alerts)
- CSS Custom Properties (for theming)
- Local Storage (for theme persistence)

**Tested Browsers:**
- Chrome 90+
- Firefox 90+
- Safari 14+
- Edge 90+

## 🛡️ Security Features

- **CORS Protection**: Configurable origin checking
- **Input Sanitization**: All user inputs are properly escaped
- **Connection Limits**: WebSocket connection management
- **SQL Injection Prevention**: Parameterized database queries
- **XSS Protection**: Content Security Policy headers