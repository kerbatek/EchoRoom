#!/bin/bash

# EchoRoom VPS Deployment Script
# Simple VPS deployment for EchoRoom chat application

set -e

# Configuration
VPS_HOST=${VPS_HOST:-"your-vps-ip"}
VPS_USER=${VPS_USER:-"root"}
VPS_PORT=${VPS_PORT:-22}
APP_NAME="echoroom"
DEPLOY_DIR="/opt/$APP_NAME"
BACKUP_DIR="/opt/$APP_NAME/backups"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging function
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}"
}

warning() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}"
}

# Check if SSH key exists
check_ssh_key() {
    if [ ! -f ~/.ssh/id_rsa ]; then
        error "SSH key not found. Please generate one with: ssh-keygen -t rsa -b 4096"
        exit 1
    fi
}

# Test SSH connection
test_ssh() {
    log "Testing SSH connection to $VPS_USER@$VPS_HOST:$VPS_PORT..."
    
    if ssh -o ConnectTimeout=10 -p $VPS_PORT $VPS_USER@$VPS_HOST "echo 'SSH connection successful'"; then
        log "✅ SSH connection successful"
    else
        error "❌ SSH connection failed"
        exit 1
    fi
}

# Setup VPS prerequisites
setup_vps() {
    log "Setting up VPS prerequisites..."
    
    ssh -p $VPS_PORT $VPS_USER@$VPS_HOST << 'EOF'
        # Update system
        apt-get update -y
        apt-get upgrade -y
        
        # Install Docker
        if ! command -v docker &> /dev/null; then
            curl -fsSL https://get.docker.com -o get-docker.sh
            sh get-docker.sh
            systemctl enable docker
            systemctl start docker
        fi
        
        # Install Docker Compose
        if ! command -v docker-compose &> /dev/null; then
            curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
            chmod +x /usr/local/bin/docker-compose
        fi
        
        # Install additional tools
        apt-get install -y curl wget git htop ufw
        
        # Setup firewall
        ufw --force enable
        ufw allow ssh
        ufw allow 80
        ufw allow 443
        ufw allow 8080
        
        # Create application directories
        mkdir -p /opt/echoroom /opt/echoroom/backups
        
        echo "✅ VPS setup completed"
EOF
    
    log "VPS prerequisites installed"
}

# Create deployment package
create_deployment_package() {
    log "Creating deployment package..."
    
    # Create temporary directory
    TEMP_DIR=$(mktemp -d)
    
    # Copy necessary files
    cp Dockerfile $TEMP_DIR/
    cp docker-compose.prod.yml $TEMP_DIR/
    cp deploy.sh $TEMP_DIR/
    cp nginx.conf $TEMP_DIR/
    cp init.sql $TEMP_DIR/
    cp -r assets $TEMP_DIR/ 2>/dev/null || true
    
    # Create VPS-specific environment file
    cat > $TEMP_DIR/.env << EOF
DB_PASSWORD=${DB_PASSWORD:-defaultpassword}
DOMAIN=${DOMAIN:-localhost}
GIN_MODE=release
EOF
    
    # Create deployment script for VPS
    cat > $TEMP_DIR/vps-deploy.sh << 'EOF'
#!/bin/bash
set -e

echo "🚀 Starting VPS deployment..."

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | xargs)
fi

# Pull latest image
docker pull ${DOCKER_USERNAME:-kerbatek}/echoroom:latest

# Stop existing services
if [ -f docker-compose.prod.yml ]; then
    docker-compose -f docker-compose.prod.yml down || true
fi

# Backup database
if docker ps -a | grep -q postgres; then
    echo "📦 Creating database backup..."
    docker exec postgres pg_dump -U postgres chat_app > "backup-$(date +%Y%m%d-%H%M%S).sql" 2>/dev/null || true
fi

# Start services
docker-compose -f docker-compose.prod.yml up -d

# Wait for services
echo "⏳ Waiting for services to start..."
sleep 30

# Health check
for i in {1..10}; do
    if curl -f http://localhost:8080/health; then
        echo "✅ Health check passed!"
        break
    else
        echo "⏳ Health check attempt $i failed, retrying..."
        sleep 5
    fi
    
    if [ $i -eq 10 ]; then
        echo "❌ Health check failed after 10 attempts"
        exit 1
    fi
done

echo "🎉 VPS deployment completed successfully!"
EOF
    
    chmod +x $TEMP_DIR/vps-deploy.sh
    
    # Create archive
    tar -czf deployment.tar.gz -C $TEMP_DIR .
    
    # Cleanup
    rm -rf $TEMP_DIR
    
    log "Deployment package created: deployment.tar.gz"
}

# Deploy to VPS
deploy_to_vps() {
    log "Deploying to VPS..."
    
    # Copy deployment package
    scp -P $VPS_PORT deployment.tar.gz $VPS_USER@$VPS_HOST:/tmp/
    
    # Execute deployment
    ssh -p $VPS_PORT $VPS_USER@$VPS_HOST << 'EOF'
        cd /opt/echoroom
        
        # Backup current deployment
        if [ -d "current" ]; then
            cp -r current backups/backup-$(date +%Y%m%d-%H%M%S)
        fi
        
        # Extract new deployment
        rm -rf new
        mkdir -p new
        cd new
        tar -xzf /tmp/deployment.tar.gz
        
        # Execute deployment
        chmod +x vps-deploy.sh
        ./vps-deploy.sh
        
        # Switch to new deployment
        cd /opt/echoroom
        rm -rf current
        mv new current
        
        # Cleanup
        rm -f /tmp/deployment.tar.gz
        
        echo "✅ VPS deployment completed!"
EOF
    
    log "Deployment successful!"
}

# Health check
health_check() {
    log "Performing health check..."
    
    if curl -f http://$VPS_HOST:8080/health; then
        log "✅ Health check passed!"
        log "🔗 Application URL: http://$VPS_HOST:8080"
    else
        error "❌ Health check failed!"
        exit 1
    fi
}

# Rollback function
rollback() {
    warning "Rolling back deployment..."
    
    ssh -p $VPS_PORT $VPS_USER@$VPS_HOST << 'EOF'
        cd /opt/echoroom
        
        # Find latest backup
        LATEST_BACKUP=$(ls -t backups/ | head -1)
        
        if [ -n "$LATEST_BACKUP" ]; then
            echo "🔄 Rolling back to: $LATEST_BACKUP"
            
            # Stop current services
            cd current
            docker-compose -f docker-compose.prod.yml down || true
            
            # Restore backup
            cd /opt/echoroom
            rm -rf current
            cp -r backups/$LATEST_BACKUP current
            
            # Start services
            cd current
            docker-compose -f docker-compose.prod.yml up -d
            
            echo "✅ Rollback completed"
        else
            echo "❌ No backup found for rollback"
        fi
EOF
}

# Show logs
show_logs() {
    log "Showing application logs..."
    ssh -p $VPS_PORT $VPS_USER@$VPS_HOST << 'EOF'
        cd /opt/echoroom/current
        docker-compose -f docker-compose.prod.yml logs -f
EOF
}

# Show status
show_status() {
    log "Showing application status..."
    ssh -p $VPS_PORT $VPS_USER@$VPS_HOST << 'EOF'
        cd /opt/echoroom/current
        docker-compose -f docker-compose.prod.yml ps
EOF
}

# Main function
main() {
    case "${1:-deploy}" in
        setup)
            check_ssh_key
            test_ssh
            setup_vps
            ;;
        deploy)
            check_ssh_key
            test_ssh
            create_deployment_package
            deploy_to_vps
            health_check
            ;;
        rollback)
            check_ssh_key
            test_ssh
            rollback
            ;;
        logs)
            check_ssh_key
            show_logs
            ;;
        status)
            check_ssh_key
            show_status
            ;;
        health)
            health_check
            ;;
        *)
            echo "Usage: $0 {setup|deploy|rollback|logs|status|health}"
            echo ""
            echo "Commands:"
            echo "  setup    - Setup VPS prerequisites (Docker, firewall, etc.)"
            echo "  deploy   - Deploy application to VPS"
            echo "  rollback - Rollback to previous deployment"
            echo "  logs     - Show application logs"
            echo "  status   - Show application status"
            echo "  health   - Perform health check"
            echo ""
            echo "Environment variables:"
            echo "  VPS_HOST - VPS IP address or hostname"
            echo "  VPS_USER - SSH username (default: root)"
            echo "  VPS_PORT - SSH port (default: 22)"
            echo "  DB_PASSWORD - Database password"
            echo "  DOMAIN - Domain name for SSL"
            exit 1
            ;;
    esac
}

# Run main function
main "$@"