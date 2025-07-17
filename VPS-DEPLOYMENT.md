# EchoRoom VPS Deployment Guide

Complete guide for deploying EchoRoom chat application to your VPS server with automated CI/CD.

## 🚀 Quick Start

### Prerequisites
- VPS with Ubuntu 20.04+ (2GB RAM, 20GB storage minimum)
- SSH access to your VPS
- Domain name (optional but recommended)
- Docker Hub account

### 1. Setup GitHub Secrets

Add these secrets to your GitHub repository settings:

| Secret Name | Description | Example |
|-------------|-------------|---------|
| `VPS_HOST` | VPS IP address or hostname | `192.168.1.100` |
| `VPS_USER` | SSH username | `root` or `ubuntu` |
| `VPS_SSH_KEY` | Private SSH key content | `-----BEGIN RSA PRIVATE KEY-----...` |
| `VPS_PORT` | SSH port (optional) | `22` |
| `DB_PASSWORD` | Database password | `your_secure_password` |
| `DOCKER_USERNAME` | Docker Hub username | `your_username` |
| `DOCKER_PASSWORD` | Docker Hub password | `your_password` |

### 2. Automated Deployment

Once secrets are configured, deployment is automatic:

1. **Push to main branch** → Triggers CI/CD pipeline
2. **Tests pass** → Builds and pushes Docker image
3. **Deploys to VPS** → Automatically via SSH
4. **Health check** → Verifies deployment
5. **Notifications** → Confirms success

## 🛠️ Manual Deployment

### Setup VPS Prerequisites

```bash
# Make script executable
chmod +x vps-deploy.sh

# Set environment variables
export VPS_HOST="your-vps-ip"
export VPS_USER="root"
export DB_PASSWORD="your_secure_password"

# Setup VPS (one-time only)
./vps-deploy.sh setup
```

### Deploy Application

```bash
# Deploy to VPS
./vps-deploy.sh deploy

# Check status
./vps-deploy.sh status

# View logs
./vps-deploy.sh logs
```

### Rollback if Needed

```bash
# Rollback to previous version
./vps-deploy.sh rollback
```

## 🔧 VPS Configuration

### 1. Initial VPS Setup

```bash
# Connect to your VPS
ssh root@your-vps-ip

# Update system
apt update && apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Install Docker Compose
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Setup firewall
ufw enable
ufw allow ssh
ufw allow 80
ufw allow 443
ufw allow 8080
```

### 2. SSL Certificate Setup (Optional)

```bash
# Install Certbot
apt install certbot python3-certbot-nginx -y

# Get SSL certificate
certbot certonly --standalone -d yourdomain.com

# Setup auto-renewal
crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

### 3. Nginx Configuration

```bash
# Install Nginx
apt install nginx -y

# Copy nginx configuration
cp nginx.conf /etc/nginx/sites-available/echoroom
ln -s /etc/nginx/sites-available/echoroom /etc/nginx/sites-enabled/

# Test configuration
nginx -t

# Restart Nginx
systemctl restart nginx
```

## 📂 Directory Structure on VPS

```
/opt/echoroom/
├── current/              # Current deployment
│   ├── docker-compose.prod.yml
│   ├── deploy.sh
│   ├── nginx.conf
│   └── init.sql
├── backups/              # Deployment backups
│   ├── backup-20240101-120000/
│   └── backup-20240101-140000/
└── logs/                 # Application logs
```

## 🔍 Monitoring and Maintenance

### Check Application Status

```bash
# Check running containers
docker ps

# Check logs
docker-compose -f /opt/echoroom/current/docker-compose.prod.yml logs -f

# Check resource usage
docker stats

# Check system resources
htop
df -h
```

### Database Management

```bash
# Create database backup
docker exec postgres pg_dump -U postgres chat_app > backup.sql

# Restore database
docker exec -i postgres psql -U postgres chat_app < backup.sql

# Connect to database
docker exec -it postgres psql -U postgres chat_app
```

### Update Application

```bash
# Pull latest image
docker pull kerbatek/echoroom:latest

# Restart services
cd /opt/echoroom/current
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d
```

## 🚨 Troubleshooting

### Common Issues

1. **SSH Connection Failed**
   ```bash
   # Check SSH service
   systemctl status ssh
   
   # Check firewall
   ufw status
   
   # Test connection
   ssh -v user@host
   ```

2. **Docker Service Not Running**
   ```bash
   # Start Docker
   systemctl start docker
   systemctl enable docker
   ```

3. **Port Already in Use**
   ```bash
   # Check what's using port 8080
   lsof -i :8080
   
   # Kill process if needed
   kill -9 PID
   ```

4. **Database Connection Issues**
   ```bash
   # Check database logs
   docker logs postgres
   
   # Check environment variables
   docker exec app env | grep DB_
   ```

### Health Check Endpoints

- **Application**: `http://your-vps-ip:8080/health`
- **Database**: Check via application logs

### Log Files

- **Application**: `docker-compose logs app`
- **Database**: `docker-compose logs postgres`
- **Nginx**: `/var/log/nginx/error.log`
- **System**: `/var/log/syslog`

## 🔐 Security Best Practices

### 1. SSH Security

```bash
# Disable password authentication
vim /etc/ssh/sshd_config
# Set: PasswordAuthentication no
systemctl restart ssh

# Use non-standard SSH port
vim /etc/ssh/sshd_config
# Set: Port 2222
systemctl restart ssh
```

### 2. Firewall Configuration

```bash
# Basic firewall rules
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow 80
ufw allow 443
ufw enable
```

### 3. Database Security

```bash
# Use strong passwords
# Enable SSL connections
# Restrict network access
# Regular backups
```

### 4. Application Security

```bash
# Keep containers updated
docker system prune -a

# Use secrets for sensitive data
# Enable HTTPS
# Implement rate limiting
```

## 📊 Performance Optimization

### 1. Resource Limits

```yaml
# In docker-compose.prod.yml
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
```

### 2. Database Optimization

```sql
-- PostgreSQL configuration
shared_preload_libraries = 'pg_stat_statements'
max_connections = 100
shared_buffers = 128MB
effective_cache_size = 512MB
```

### 3. Nginx Optimization

```nginx
# In nginx.conf
worker_processes auto;
worker_connections 1024;

gzip on;
gzip_types text/plain text/css application/json application/javascript;
```

## 🔄 Backup Strategy

### Automated Backups

```bash
# Create backup script
cat > /opt/echoroom/backup.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/opt/echoroom/backups"
DATE=$(date +%Y%m%d_%H%M%S)

# Database backup
docker exec postgres pg_dump -U postgres chat_app > "$BACKUP_DIR/db_$DATE.sql"

# Application backup
tar -czf "$BACKUP_DIR/app_$DATE.tar.gz" /opt/echoroom/current

# Cleanup old backups (keep last 7 days)
find "$BACKUP_DIR" -name "*.sql" -mtime +7 -delete
find "$BACKUP_DIR" -name "*.tar.gz" -mtime +7 -delete
EOF

chmod +x /opt/echoroom/backup.sh

# Setup daily backup cron
crontab -e
# Add: 0 2 * * * /opt/echoroom/backup.sh
```

## 📞 Support and Maintenance

### Regular Maintenance Tasks

1. **System Updates** (Monthly)
   ```bash
   apt update && apt upgrade -y
   ```

2. **Docker Cleanup** (Weekly)
   ```bash
   docker system prune -f
   ```

3. **Log Rotation** (Automatic)
   ```bash
   # Setup logrotate for application logs
   ```

4. **SSL Certificate Renewal** (Automatic)
   ```bash
   # Certbot handles this automatically
   ```

### Emergency Procedures

1. **Application Down**
   - Check container status
   - Review logs
   - Restart services
   - Contact support if needed

2. **Database Issues**
   - Check disk space
   - Review database logs
   - Restore from backup if needed

3. **SSL Certificate Expired**
   - Renew certificate manually
   - Restart nginx
   - Update DNS if needed

## 🎯 Next Steps

1. **Setup monitoring** (Prometheus, Grafana)
2. **Configure log aggregation** (ELK stack)
3. **Implement alerting** (email, Slack)
4. **Setup load balancing** (for high traffic)
5. **Add CDN** (for static assets)

For additional help, check:
- [Main Deployment Guide](DEPLOYMENT.md)
- [GitHub Issues](https://github.com/kerbatek/EchoRoom/issues)
- [Docker Documentation](https://docs.docker.com/)
- [Nginx Documentation](https://nginx.org/en/docs/)