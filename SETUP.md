# EchoRoom VPS Setup Guide

Complete step-by-step guide to set up automated VPS deployment for EchoRoom.

## 📋 Prerequisites

- VPS server (Ubuntu 20.04+, 2GB RAM, 20GB storage)
- GitHub account with repository access
- Docker Hub account
- Domain name (optional)
- SSH key pair

## 🔧 Step 1: Prepare Your VPS

### 1.1 Connect to Your VPS

```bash
# Connect via SSH
ssh root@your-vps-ip

# Or if using ubuntu user
ssh ubuntu@your-vps-ip
```

### 1.2 Update System

```bash
# Update package lists
sudo apt update && sudo apt upgrade -y

# Install essential packages
sudo apt install -y curl wget git htop ufw
```

### 1.3 Install Docker

```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Start and enable Docker
sudo systemctl start docker
sudo systemctl enable docker

# Add user to docker group (if not using root)
sudo usermod -aG docker $USER
```

### 1.4 Install Docker Compose

```bash
# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Verify installation
docker --version
docker-compose --version
```

### 1.5 Setup Firewall

```bash
# Enable firewall
sudo ufw --force enable

# Allow SSH (change port if needed)
sudo ufw allow 22

# Allow HTTP and HTTPS
sudo ufw allow 80
sudo ufw allow 443

# Allow application port
sudo ufw allow 8080

# Check status
sudo ufw status
```

## 🔑 Step 2: Setup SSH Keys

### 2.1 Generate SSH Key (if you don't have one)

```bash
# On your local machine
ssh-keygen -t rsa -b 4096 -C "your-email@example.com"

# Press Enter to accept default location
# Set a passphrase (optional)
```

### 2.2 Copy SSH Key to VPS

```bash
# Copy public key to VPS
ssh-copy-id root@your-vps-ip

# Or manually copy
cat ~/.ssh/id_rsa.pub | ssh root@your-vps-ip "mkdir -p ~/.ssh && cat >> ~/.ssh/authorized_keys"
```

### 2.3 Test SSH Connection

```bash
# Test passwordless SSH
ssh root@your-vps-ip

# Should connect without password
```

## 🐳 Step 3: Setup Docker Hub

### 3.1 Create Docker Hub Account

1. Go to [Docker Hub](https://hub.docker.com/)
2. Create account or sign in
3. Create a repository named `echoroom` (lowercase)

### 3.2 Get Docker Hub Credentials

- **Username**: Your Docker Hub username
- **Password**: Your Docker Hub password or access token

## 🔐 Step 4: Configure GitHub Secrets

### 4.1 Access GitHub Repository Settings

1. Go to your GitHub repository
2. Click **Settings** tab
3. Click **Secrets and variables** → **Actions**
4. Click **New repository secret**

### 4.2 Add Required Secrets

Add these 6 secrets one by one:

#### VPS_HOST
```
Name: VPS_HOST
Value: 192.168.1.100  # Your VPS IP address
```

#### VPS_USER
```
Name: VPS_USER
Value: root  # or ubuntu
```

#### VPS_SSH_KEY
```
Name: VPS_SSH_KEY
Value: # Content of your private SSH key
```

To get your private SSH key:
```bash
# On your local machine
cat ~/.ssh/id_rsa
```
Copy the ENTIRE output including:
```
-----BEGIN RSA PRIVATE KEY-----
...your key content...
-----END RSA PRIVATE KEY-----
```

#### DB_PASSWORD
```
Name: DB_PASSWORD
Value: your_secure_database_password_123
```

#### DOCKER_USERNAME
```
Name: DOCKER_USERNAME
Value: your-dockerhub-username
```

#### DOCKER_PASSWORD
```
Name: DOCKER_PASSWORD
Value: your-dockerhub-password
```

## 🚀 Step 5: Test Initial Setup

### 5.1 Setup VPS Environment

```bash
# Clone repository to your local machine
git clone https://github.com/kerbatek/EchoRoom.git
cd EchoRoom

# Run VPS setup (one-time only)
export VPS_HOST="your-vps-ip"
export VPS_USER="root"
export DB_PASSWORD="your_secure_password"

./vps-deploy.sh setup
```

### 5.2 Test Manual Deployment

```bash
# Test manual deployment
./vps-deploy.sh deploy

# Check status
./vps-deploy.sh status

# View logs
./vps-deploy.sh logs
```

## 📱 Step 6: Test Automated Deployment

### 6.1 Push to Main Branch

```bash
# Make a small change
echo "# Test deployment" >> README.md

# Commit and push
git add .
git commit -m "test: trigger automated deployment"
git push origin main
```

### 6.2 Monitor Deployment

1. Go to your GitHub repository
2. Click **Actions** tab
3. Watch the deployment workflow
4. Check each step completion

### 6.3 Verify Deployment

```bash
# Check if app is running
curl http://your-vps-ip:8080/health

# Should return: {"status":"ok"}
```

## 🌐 Step 7: Access Your Application

### 7.1 Open in Browser

```
http://your-vps-ip:8080
```

### 7.2 Test Functionality

1. Enter a username
2. Send messages
3. Switch between channels
4. Create new channels

## 🔍 Step 8: Troubleshooting

### 8.1 Common Issues

#### GitHub Actions Fails
```bash
# Check GitHub Actions logs
# Go to Actions tab → Click on failed run → Check logs

# Common fixes:
# 1. Verify all 6 secrets are set correctly
# 2. Check SSH key format (include BEGIN/END lines)
# 3. Verify VPS connectivity
```

#### SSH Connection Issues
```bash
# Test SSH connection manually
ssh -i ~/.ssh/id_rsa root@your-vps-ip

# Check SSH key permissions
chmod 600 ~/.ssh/id_rsa
chmod 644 ~/.ssh/id_rsa.pub
```

#### Docker Issues
```bash
# Check Docker status on VPS
ssh root@your-vps-ip "docker ps"

# Restart Docker if needed
ssh root@your-vps-ip "sudo systemctl restart docker"
```

#### Application Not Starting
```bash
# Check application logs
./vps-deploy.sh logs

# Check container status
ssh root@your-vps-ip "cd /opt/echoroom/current && docker-compose ps"
```

### 8.2 Debug Commands

```bash
# Check VPS status
./vps-deploy.sh status

# View detailed logs
./vps-deploy.sh logs

# Test health endpoint
curl -v http://your-vps-ip:8080/health

# Check firewall
ssh root@your-vps-ip "sudo ufw status"

# Check running processes
ssh root@your-vps-ip "ps aux | grep docker"
```

## 🔄 Step 9: Maintenance

### 9.1 Regular Updates

```bash
# Updates happen automatically on push to main
git push origin main
```

### 9.2 Manual Operations

```bash
# Manual deployment
./vps-deploy.sh deploy

# Check status
./vps-deploy.sh status

# View logs
./vps-deploy.sh logs

# Rollback if needed
./vps-deploy.sh rollback
```

### 9.3 VPS Maintenance

```bash
# Update VPS system (monthly)
ssh root@your-vps-ip "sudo apt update && sudo apt upgrade -y"

# Clean up Docker (weekly)
ssh root@your-vps-ip "docker system prune -f"

# Check disk space
ssh root@your-vps-ip "df -h"
```

## 📊 Step 10: Monitoring

### 10.1 Health Checks

```bash
# Application health
curl http://your-vps-ip:8080/health

# Database health (via SSH)
ssh root@your-vps-ip "docker exec postgres pg_isready -U postgres"
```

### 10.2 Log Monitoring

```bash
# Real-time logs
./vps-deploy.sh logs

# Application logs only
ssh root@your-vps-ip "cd /opt/echoroom/current && docker-compose logs -f app"

# Database logs
ssh root@your-vps-ip "cd /opt/echoroom/current && docker-compose logs -f postgres"
```

## 🎯 Next Steps

1. **Custom Domain**: Point your domain to VPS IP
2. **SSL Certificate**: Setup Let's Encrypt for HTTPS
3. **Monitoring**: Add application monitoring
4. **Backups**: Setup automated database backups
5. **Scaling**: Add load balancer for high traffic

## 📞 Support

If you encounter issues:

1. **Check this guide** for common solutions
2. **Review GitHub Actions logs** for deployment errors
3. **Test manual commands** using `./vps-deploy.sh`
4. **Check VPS logs** for runtime issues
5. **Create GitHub issue** for persistent problems

## ✅ Checklist

- [ ] VPS prepared with Docker
- [ ] SSH keys configured
- [ ] GitHub secrets added (all 6)
- [ ] Docker Hub account setup
- [ ] Initial deployment tested
- [ ] Automated deployment working
- [ ] Application accessible
- [ ] Health checks passing

**Congratulations!** 🎉 Your EchoRoom deployment is ready!