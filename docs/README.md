# Real-time Chat Application

A WebSocket-based chat application with PostgreSQL persistence, supporting both ephemeral and persistent channels.

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
go run *.go
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

## Features

- **Ephemeral Channels** ⚡: Temporary channels that disappear when empty
- **Persistent Channels** 💾: Permanent channels with message history
- **Random Usernames**: Automatically assigned funny usernames
- **Real-time Sync**: All clients see channel changes instantly
- **Message History**: Last 50 messages loaded for persistent channels
- **Dark/Light Theme**: Toggle between themes with localStorage persistence

## Architecture

### Backend (Modular Go)
- `main.go` - Application entry point and routing
- `types.go` - Type definitions (Hub, Channel, Client, Message)
- `hub.go` - Central hub managing clients and channels
- `channel.go` - Channel broadcast and lifecycle management
- `client.go` - WebSocket client connection handling
- `websocket.go` - WebSocket endpoint and upgrade logic
- `database.go` - PostgreSQL operations and schema management

### Frontend (Single Page)
- `index.html` - Complete frontend with embedded CSS/JavaScript
- WebSocket connection with automatic reconnection
- Real-time channel management and message display
- Responsive design with sidebar and chat area

### Database Schema
**Channels Table:**
- `name` (VARCHAR PRIMARY KEY) - Channel identifier
- `type` (VARCHAR) - "ephemeral" or "persistent"
- `created_at` (TIMESTAMP) - Creation time

**Messages Table:**
- `id` (SERIAL PRIMARY KEY) - Message identifier
- `channel_name` (VARCHAR) - References channels.name
- `username` (VARCHAR) - Message author
- `content` (TEXT) - Message content
- `timestamp` (TIMESTAMP) - Message time

## Commands

### Development
```bash
go run *.go                 # Start development server
go build -o chat-app *.go   # Build executable
go mod tidy                 # Clean up dependencies
```

### Testing
```bash
go test ./...               # Run all tests
go test -v ./...           # Verbose test output
go test -cover ./...       # Test coverage report
```

## Environment Variables

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