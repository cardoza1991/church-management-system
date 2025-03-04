# Church Management System Documentation

## Project Overview
This system helps track contacts ("fruit") through their spiritual journey, from initial contact to becoming gospel workers.

## Architecture
- Microservices architecture with Golang
- MariaDB database
- Kubernetes orchestration
- Wappler frontend

## Services

### User Service
Handles authentication and user management.

### Contact Service
Manages contacts and their spiritual journey.

### Study Service
Tracks Bible studies and sessions.

### Reservation Service
Handles room reservations.

## Development Setup
1. Install Minikube and kubectl
2. Start Minikube: `minikube start`
3. Apply Kubernetes configs: `kubectl apply -f infrastructure/k8s/`
4. Run database migrations
5. Start developing services

## API Documentation
[To be added]
