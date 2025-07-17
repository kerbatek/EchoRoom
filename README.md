# EchoRoom 🎯

A real-time chat application built with Go, WebSocket, and PostgreSQL.

## Features ✨

- **Real-time messaging** with WebSocket connections
- **Multiple channels** (persistent and ephemeral)
- **Channel switching** without page reload
- **Message history** for persistent channels
- **Race condition free** with proper synchronization
- **Docker containerized** for easy deployment

## Quick Start 🚀

### Local Development

```bash
# Clone repository
git clone https://github.com/kerbatek/EchoRoom.git
cd EchoRoom

# Install dependencies
go mod download

# Setup database
# Configure PostgreSQL and update connection settings

# Run application
go run .
```

### Docker

```bash
# Build Docker image
docker build -t echoroom .

# Run with Docker
docker run -p 8080:8080 echoroom
```

## Architecture 🏗️

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   GitHub Repo   │───▶│  GitHub Actions │───▶│  Docker Registry│
│                 │    │                 │    │                 │
│ - Push to main  │    │ - Run tests     │    │ - Store image   │
│ - Trigger CI/CD │    │ - Build image   │    │ - Tagged builds │
│                 │    │ - Push to registry│    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Tech Stack 🛠️

- **Backend**: Go (Gin framework)
- **Database**: PostgreSQL
- **WebSocket**: Gorilla WebSocket
- **Containerization**: Docker
- **CI/CD**: GitHub Actions

## Development 💻

### Testing

```bash
# Run tests
go test ./...

# Run tests with race detection
go test -race ./...

# Run specific test
go test -v ./... -run TestHubClientRegistration
```

### Building

```bash
# Build for current platform
go build -v ./...

# Build for Linux
GOOS=linux GOARCH=amd64 go build -v ./...
```

## Usage 📱

1. **Access** the application at `http://localhost:8080`
2. **Enter username** and start chatting
3. **Switch channels** using the channel list
4. **Create channels** (persistent/ephemeral)
5. **View history** in persistent channels

## Configuration ⚙️

### Environment Variables

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=chat_app
DB_SSLMODE=disable
GIN_MODE=release
```

### Database Setup

The application requires PostgreSQL. Create a database and update the connection settings in your environment variables.

## Monitoring 📊

### Health Check

```bash
curl http://localhost:8080/health
```

### Logs

```bash
# View application logs
docker logs container_name

# Follow logs
docker logs -f container_name
```

## Security 🔒

- **Database Password** protection
- **Docker Container** isolation
- **Input validation** for all user inputs
- **WebSocket** secure connections

## Contributing 🤝

1. Fork the repository
2. Create feature branch
3. Make changes
4. Run tests
5. Submit pull request

## License 📄

This project is licensed under the MIT License.

## Support 💬

- **Issues**: [GitHub Issues](https://github.com/kerbatek/EchoRoom/issues)
- **GitHub Actions**: Check workflow logs for build issues

---

**Happy Chatting!** 🎉