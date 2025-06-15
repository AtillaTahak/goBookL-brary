#!/bin/bash

# Go Book Library Monitoring Setup Script
# This script helps set up and manage the monitoring stack (Prometheus, Grafana, Alertmanager)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}======================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}======================================${NC}"
}

# Check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker and try again."
        exit 1
    fi
    print_status "Docker is running"
}

# Check if docker-compose is available
check_docker_compose() {
    if ! command -v docker-compose &> /dev/null; then
        print_error "docker-compose is not installed. Please install it and try again."
        exit 1
    fi
    print_status "docker-compose is available"
}

# Start the monitoring stack
start_monitoring() {
    print_header "Starting Monitoring Stack"

    cd "$(dirname "$0")"

    print_status "Starting Prometheus, Grafana, and Alertmanager..."
    docker-compose up -d prometheus grafana alertmanager redis-exporter postgres-exporter node-exporter

    print_status "Waiting for services to be ready..."
    sleep 30

    # Check if services are running
    if docker-compose ps | grep -q "prometheus.*Up"; then
        print_status "✅ Prometheus is running at http://localhost:9090"
    else
        print_error "❌ Prometheus failed to start"
    fi

    if docker-compose ps | grep -q "grafana.*Up"; then
        print_status "✅ Grafana is running at http://localhost:3001 (admin/admin)"
    else
        print_error "❌ Grafana failed to start"
    fi

    if docker-compose ps | grep -q "alertmanager.*Up"; then
        print_status "✅ Alertmanager is running at http://localhost:9093"
    else
        print_error "❌ Alertmanager failed to start"
    fi
}

# Start the full application stack
start_full() {
    print_header "Starting Full Application Stack"

    cd "$(dirname "$0")"

    print_status "Building and starting all services..."
    docker-compose up -d --build

    print_status "Waiting for all services to be ready..."
    sleep 60

    print_status "Application stack started!"
    print_status "🚀 Frontend: http://localhost:3000"
    print_status "🔧 Backend API: http://localhost:8080"
    print_status "📊 Prometheus: http://localhost:9090"
    print_status "📈 Grafana: http://localhost:3001"
    print_status "🚨 Alertmanager: http://localhost:9093"
    print_status "💾 PostgreSQL: localhost:5432"
    print_status "⚡ Redis: localhost:6379"
}

# Stop all services
stop_services() {
    print_header "Stopping All Services"

    cd "$(dirname "$0")"
    docker-compose down

    print_status "All services stopped"
}

# Show service status
show_status() {
    print_header "Service Status"

    cd "$(dirname "$0")"
    docker-compose ps
}

# Show logs for a specific service
show_logs() {
    if [ -z "$1" ]; then
        print_error "Please specify a service name"
        print_status "Available services: backend, frontend, postgres, redis, prometheus, grafana, alertmanager"
        return 1
    fi

    cd "$(dirname "$0")"
    docker-compose logs -f "$1"
}

# Clean up everything (including volumes)
cleanup() {
    print_header "Cleaning Up Everything"

    read -p "This will remove all data including databases. Are you sure? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        cd "$(dirname "$0")"
        docker-compose down -v --remove-orphans
        docker system prune -f
        print_status "Cleanup completed"
    else
        print_status "Cleanup cancelled"
    fi
}

# Health check for all services
health_check() {
    print_header "Health Check"

    cd "$(dirname "$0")"

    # Check Backend
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        print_status "✅ Backend API is healthy"
    else
        print_warning "❌ Backend API is not responding"
    fi

    # Check Prometheus
    if curl -f http://localhost:9090/-/healthy > /dev/null 2>&1; then
        print_status "✅ Prometheus is healthy"
    else
        print_warning "❌ Prometheus is not responding"
    fi

    # Check Grafana
    if curl -f http://localhost:3001/api/health > /dev/null 2>&1; then
        print_status "✅ Grafana is healthy"
    else
        print_warning "❌ Grafana is not responding"
    fi

    # Check Alertmanager
    if curl -f http://localhost:9093/-/healthy > /dev/null 2>&1; then
        print_status "✅ Alertmanager is healthy"
    else
        print_warning "❌ Alertmanager is not responding"
    fi
}

# Show help
show_help() {
    echo "Go Book Library Monitoring Setup Script"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  start-monitoring  Start only monitoring stack (Prometheus, Grafana, Alertmanager)"
    echo "  start-full       Start the complete application stack"
    echo "  stop             Stop all services"
    echo "  status           Show service status"
    echo "  logs [service]   Show logs for a specific service"
    echo "  health           Check health of all services"
    echo "  cleanup          Stop and remove all containers and volumes"
    echo "  help             Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 start-full"
    echo "  $0 logs backend"
    echo "  $0 health"
}

# Main script logic
main() {
    case "${1:-help}" in
        "start-monitoring")
            check_docker
            check_docker_compose
            start_monitoring
            ;;
        "start-full")
            check_docker
            check_docker_compose
            start_full
            ;;
        "stop")
            stop_services
            ;;
        "status")
            show_status
            ;;
        "logs")
            show_logs "$2"
            ;;
        "health")
            health_check
            ;;
        "cleanup")
            cleanup
            ;;
        "help"|*)
            show_help
            ;;
    esac
}

main "$@"
