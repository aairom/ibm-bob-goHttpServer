# Go HTTP Server

A production-ready HTTP server written in Go with Docker and Kubernetes support for local deployment on minikube.

## Features

- ✅ RESTful API with multiple endpoints
- ✅ Health check endpoint for monitoring
- ✅ Request logging middleware
- ✅ CORS support
- ✅ Graceful shutdown
- ✅ Docker containerization with multi-stage builds
- ✅ Kubernetes deployment manifests
- ✅ Security best practices (non-root user, minimal image)

## Project Structure

```
bob-project1/
├── main.go                 # Main HTTP server application
├── go.mod                  # Go module dependencies
├── Dockerfile              # Docker image configuration
├── .dockerignore          # Files to exclude from Docker build
├── k8s/
│   ├── deployment.yaml    # Kubernetes Deployment manifest
│   └── service.yaml       # Kubernetes Service manifest (NodePort)
├── Makefile               # Build and deployment automation
└── README.md              # This file
```

## Prerequisites

- **Go 1.21+** - [Install Go](https://golang.org/doc/install)
- **Docker** - [Install Docker](https://docs.docker.com/get-docker/)
- **Minikube** - [Install Minikube](https://minikube.sigs.k8s.io/docs/start/)
- **kubectl** - [Install kubectl](https://kubernetes.io/docs/tasks/tools/)
- **make** (optional) - For using Makefile commands

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Welcome message and available endpoints |
| GET | `/health` | Health check (returns status and uptime) |
| GET | `/api/info` | Server information (version, hostname, timestamp) |
| GET | `/api/echo?message=<text>` | Echo endpoint that returns the message |
| POST | `/api/data` | Demo POST endpoint (accepts and returns JSON) |

## Quick Start

### 1. Run Locally (Without Docker)

```bash
# Run the server
go run main.go

# Or use Makefile
make run
```

The server will start on `http://localhost:8080`

### 2. Test the Endpoints

```bash
# Test root endpoint
curl http://localhost:8080/

# Test health check
curl http://localhost:8080/health

# Test info endpoint
curl http://localhost:8080/api/info

# Test echo endpoint
curl "http://localhost:8080/api/echo?message=Hello"

# Test data endpoint (POST)
curl -X POST http://localhost:8080/api/data \
  -H "Content-Type: application/json" \
  -d '{"name":"test","value":"data"}'

# Or use Makefile to test all endpoints
make quick-test
```

## Docker Deployment

### Build Docker Image

```bash
# Build the image
docker build -t go-http-server:latest .

# Or use Makefile
make docker-build
```

### Run Docker Container

```bash
# Run the container
docker run -p 8080:8080 go-http-server:latest
docker run -p 8081:8080 go-http-server:latest

# Or use Makefile
make docker-run
```

Access the server at `http://localhost:8080`

## Kubernetes Deployment (Minikube)

### Step 1: Start Minikube

```bash
# Start minikube cluster
minikube start

# Verify minikube is running
minikube status
```

### Step 2: Build Image in Minikube

```bash
# Configure Docker to use minikube's Docker daemon
eval $(minikube docker-env)

# Build the image in minikube
docker build -t go-http-server:latest .

# Or use Makefile
make minikube-build
```

### Step 3: Deploy to Kubernetes

```bash
# Apply Kubernetes manifests
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml

# Or use Makefile
make k8s-deploy
```

### Step 4: Access the Service

```bash
# Get the service URL
minikube service go-http-server --url

# Or open in browser
minikube service go-http-server

# Or use Makefile
make k8s-url      # Get URL
make k8s-open     # Open in browser
```

The service will be accessible at `http://<minikube-ip>:30080`

### Step 5: Verify Deployment

```bash
# Check deployment status
kubectl get deployments
kubectl get pods
kubectl get services

# Or use Makefile
make k8s-status

# View logs
kubectl logs -l app=go-http-server

# Or use Makefile
make k8s-logs
```

## Makefile Commands

The project includes a comprehensive Makefile for common operations:

```bash
# Show all available commands
make help

# Local development
make run              # Run the application locally
make build            # Build the Go binary
make test             # Run tests
make quick-test       # Test all endpoints

# Docker operations
make docker-build     # Build Docker image
make docker-run       # Run Docker container
make docker-stop      # Stop Docker container

# Minikube operations
make minikube-setup   # Setup minikube environment
make minikube-build   # Build image in minikube

# Kubernetes operations
make k8s-deploy       # Deploy to Kubernetes
make k8s-delete       # Delete Kubernetes resources
make k8s-status       # Check deployment status
make k8s-logs         # Show logs from pods
make k8s-url          # Get service URL
make k8s-open         # Open service in browser
make k8s-restart      # Restart deployment

# Combined operations
make full-deploy      # Build in minikube and deploy

# Cleanup
make clean            # Clean up build artifacts
```

## Full Deployment Workflow

For a complete deployment from scratch:

```bash
# 1. Start minikube
minikube start

# 2. Build and deploy (one command)
make full-deploy

# 3. Get the service URL
make k8s-url

# 4. Test the endpoints
MINIKUBE_IP=$(minikube ip)
curl http://$MINIKUBE_IP:30080/health
```

## Configuration

### Environment Variables

- `PORT` - Server port (default: 8080)

### Kubernetes Configuration

- **Replicas**: 2 (for high availability)
- **Resources**:
  - CPU Request: 50m
  - CPU Limit: 100m
  - Memory Request: 64Mi
  - Memory Limit: 128Mi
- **Service Type**: NodePort
- **NodePort**: 30080

## Monitoring and Debugging

### View Logs

```bash
# View logs from all pods
kubectl logs -l app=go-http-server --tail=50

# Follow logs in real-time
kubectl logs -l app=go-http-server -f

# Or use Makefile
make k8s-logs
```

### Check Pod Status

```bash
# Get pod details
kubectl describe pod -l app=go-http-server

# Get events
kubectl get events --sort-by=.metadata.creationTimestamp
```

### Access Pod Shell

```bash
# Get pod name
POD_NAME=$(kubectl get pods -l app=go-http-server -o jsonpath='{.items[0].metadata.name}')

# Execute shell in pod
kubectl exec -it $POD_NAME -- /bin/sh
```

## Troubleshooting

### Issue: Pods not starting

```bash
# Check pod status
kubectl get pods -l app=go-http-server

# Check pod events
kubectl describe pod -l app=go-http-server

# Check logs
kubectl logs -l app=go-http-server
```

### Issue: Image not found in minikube

Make sure you built the image in minikube's Docker daemon:

```bash
# Configure Docker to use minikube
eval $(minikube docker-env)

# Rebuild the image
docker build -t go-http-server:latest .
```

### Issue: Service not accessible

```bash
# Check service status
kubectl get svc go-http-server

# Get service URL
minikube service go-http-server --url

# Check if minikube is running
minikube status
```

## Cleanup

### Delete Kubernetes Resources

```bash
# Delete all resources
kubectl delete -f k8s/

# Or use Makefile
make k8s-delete
```

### Stop Minikube

```bash
# Stop minikube
minikube stop

# Delete minikube cluster
minikube delete
```

## Security Features

- ✅ Non-root user in Docker container
- ✅ Read-only root filesystem
- ✅ Dropped all capabilities
- ✅ Security context in Kubernetes
- ✅ Resource limits
- ✅ Health checks (liveness and readiness probes)

## Performance

- Multi-stage Docker build for minimal image size (~15MB)
- Efficient Go binary with static linking
- Request/response timeouts configured
- Graceful shutdown handling

## License

MIT License

## Author

Bob - Software Engineer

---

**Happy Coding! 🚀**