# EchoRoom 🎯

A real-time chat application built with Go, WebSocket, and PostgreSQL with automated VPS deployment.

## Features ✨

- **Real-time messaging** with WebSocket connections
- **Multiple channels** (persistent and ephemeral)
- **Channel switching** without page reload
- **Message history** for persistent channels
- **Automatic deployment** to VPS
- **Race condition free** with proper synchronization
- **Docker containerized** for easy deployment

## Quick Start 🚀

### 📖 Complete Setup Guide

👉 **[Follow the detailed setup guide](SETUP.md)** for step-by-step instructions.

### ⚡ Quick Summary

1. **Prepare VPS** - Install Docker, setup firewall
2. **Configure SSH** - Generate and copy SSH keys
3. **Add GitHub Secrets** - 6 required secrets
4. **Deploy** - Push to main branch for automatic deployment

```bash
# After setup, deployment is automatic
git push origin main
```

### 🔧 Manual Deployment (Optional)

```bash
# Setup VPS prerequisites (one-time)
./vps-deploy.sh setup

# Deploy application
./vps-deploy.sh deploy

# Check status
./vps-deploy.sh status
```

## Architecture 🏗️

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   GitHub Repo   │───▶│  GitHub Actions │───▶│      VPS        │
│                 │    │                 │    │                 │
│ - Push to main  │    │ - Run tests     │    │ - Pull image    │
│ - Trigger CI/CD │    │ - Build image   │    │ - Deploy app    │
│                 │    │ - Deploy to VPS │    │ - Health check  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Tech Stack 🛠️

- **Backend**: Go (Gin framework)
- **Database**: PostgreSQL
- **WebSocket**: Gorilla WebSocket
- **Containerization**: Docker
- **Deployment**: GitHub Actions
- **Infrastructure**: VPS with Docker Compose

## Development 💻

### Local Setup

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

### Testing

```bash
# Run tests
go test ./...

# Run tests with race detection
go test -race ./...

# Run specific test
go test -v ./... -run TestHubClientRegistration
```

## Deployment 🚀

### Automated Deployment (Recommended)

1. **Configure GitHub secrets** (see Quick Start)
2. **Push to main branch** - automatic deployment
3. **Monitor** in GitHub Actions tab
4. **Access** your app at `http://your-vps-ip:8080`

### Manual Deployment

See [VPS-DEPLOYMENT.md](VPS-DEPLOYMENT.md) for detailed instructions.

## Usage 📱

1. **Access** the application at your VPS IP
2. **Enter username** and start chatting
3. **Switch channels** using the channel list
4. **Create channels** (persistent/ephemeral)
5. **View history** in persistent channels

## Configuration ⚙️

### Environment Variables

```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_secure_password
DB_NAME=chat_app
DB_SSLMODE=require
GIN_MODE=release
```

### VPS Requirements

- **OS**: Ubuntu 20.04+
- **RAM**: 2GB minimum
- **Storage**: 20GB minimum
- **Docker**: Required
- **Ports**: 22 (SSH), 80, 443, 8080

## Monitoring 📊

### Health Check

```bash
curl http://your-vps-ip:8080/health
```

### Logs

```bash
# Application logs
./vps-deploy.sh logs

# Container logs
docker-compose logs -f
```

### Status

```bash
# Check deployment status
./vps-deploy.sh status

# Check containers
docker ps
```

## Security 🔒

- **SSH Key Authentication** for VPS access
- **Database Password** protection
- **Docker Container** isolation
- **Firewall Rules** for VPS
- **SSL/TLS** ready configuration

## Troubleshooting 🔧

### Common Issues

1. **Deployment fails**
   - Check GitHub secrets configuration
   - Verify VPS SSH access
   - Review GitHub Actions logs

2. **Health check fails**
   - Check if containers are running
   - Verify database connection
   - Check application logs

3. **SSH connection issues**
   - Verify SSH key format
   - Check VPS firewall settings
   - Confirm SSH service is running

### Commands

```bash
# Check deployment status
./vps-deploy.sh status

# View logs
./vps-deploy.sh logs

# Health check
./vps-deploy.sh health

# Rollback deployment
./vps-deploy.sh rollback
```

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
- **Documentation**: [VPS-DEPLOYMENT.md](VPS-DEPLOYMENT.md)
- **GitHub Actions**: Check workflow logs for deployment issues

---

**Happy Chatting!** 🎉