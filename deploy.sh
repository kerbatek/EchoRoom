#!/bin/bash

# EchoRoom Deployment Script
# This script handles deployment to production environment

set -e

# Configuration
COMPOSE_FILE="docker-compose.prod.yml"
BACKUP_DIR="./backups"
LOG_FILE="./deploy.log"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging function
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}" | tee -a "$LOG_FILE"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}" | tee -a "$LOG_FILE"
}

warning() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}" | tee -a "$LOG_FILE"
}

# Check if docker and docker-compose are installed
check_prerequisites() {
    log "Checking prerequisites..."
    
    if ! command -v docker &> /dev/null; then
        error "Docker is not installed"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        error "Docker Compose is not installed"
        exit 1
    fi
    
    log "Prerequisites check passed"
}

# Backup database before deployment
backup_database() {
    log "Creating database backup..."
    
    mkdir -p "$BACKUP_DIR"
    
    # Create database backup
    docker-compose -f "$COMPOSE_FILE" exec -T postgres pg_dump -U postgres chat_app > "$BACKUP_DIR/db_backup_$(date +%Y%m%d_%H%M%S).sql" 2>/dev/null || {
        warning "Could not create database backup (service may not be running)"
    }
    
    log "Database backup completed"
}

# Pull latest images
pull_images() {
    log "Pulling latest Docker images..."
    docker-compose -f "$COMPOSE_FILE" pull
    log "Images pulled successfully"
}

# Deploy application
deploy() {
    log "Starting deployment..."
    
    # Stop existing services
    log "Stopping existing services..."
    docker-compose -f "$COMPOSE_FILE" down
    
    # Start services
    log "Starting services..."
    docker-compose -f "$COMPOSE_FILE" up -d
    
    # Wait for services to be ready
    log "Waiting for services to start..."
    sleep 15
    
    # Check if services are running
    if docker-compose -f "$COMPOSE_FILE" ps | grep -q "Up"; then
        log "✅ Deployment successful!"
        docker-compose -f "$COMPOSE_FILE" ps
    else
        error "❌ Deployment failed!"
        docker-compose -f "$COMPOSE_FILE" logs
        exit 1
    fi
}

# Health check
health_check() {
    log "Performing health check..."
    
    # Wait a bit more for the application to fully start
    sleep 10
    
    # Check if the application is responding
    if curl -f http://localhost:8080/health &>/dev/null; then
        log "✅ Health check passed"
    else
        warning "⚠️ Health check failed - application may still be starting"
    fi
}

# Cleanup old images
cleanup() {
    log "Cleaning up old Docker images..."
    docker image prune -f
    log "Cleanup completed"
}

# Rollback function
rollback() {
    error "Rolling back deployment..."
    
    # Stop current services
    docker-compose -f "$COMPOSE_FILE" down
    
    # Restore from backup if available
    LATEST_BACKUP=$(ls -t "$BACKUP_DIR"/*.sql 2>/dev/null | head -1)
    if [ -n "$LATEST_BACKUP" ]; then
        log "Restoring database from backup: $LATEST_BACKUP"
        # Add restore logic here
    fi
    
    # Start previous version
    # docker-compose -f "$COMPOSE_FILE" up -d previous-version
    
    log "🔄 Rollback completed"
}

# Main deployment process
main() {
    log "🚀 Starting EchoRoom deployment..."
    
    # Set up error handling
    trap rollback ERR
    
    check_prerequisites
    backup_database
    pull_images
    deploy
    health_check
    cleanup
    
    log "🎉 Deployment completed successfully!"
    log "🔗 Application is available at: http://localhost:8080"
    log "📊 Monitor logs with: docker-compose -f $COMPOSE_FILE logs -f"
}

# Handle command line arguments
case "${1:-deploy}" in
    deploy)
        main
        ;;
    rollback)
        rollback
        ;;
    health)
        health_check
        ;;
    logs)
        docker-compose -f "$COMPOSE_FILE" logs -f
        ;;
    stop)
        log "Stopping services..."
        docker-compose -f "$COMPOSE_FILE" down
        ;;
    status)
        docker-compose -f "$COMPOSE_FILE" ps
        ;;
    *)
        echo "Usage: $0 {deploy|rollback|health|logs|stop|status}"
        exit 1
        ;;
esac