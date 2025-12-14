.PHONY: help run build docker-build docker-run k8s-deploy k8s-delete k8s-status clean test

# Variables
APP_NAME=go-http-server
IMAGE_NAME=go-http-server
IMAGE_TAG=latest
PORT=8080

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

run: ## Run the application locally
	@echo "Starting Go HTTP server on port $(PORT)..."
	go run main.go

build: ## Build the Go binary
	@echo "Building Go binary..."
	go build -o $(APP_NAME) main.go
	@echo "Binary created: $(APP_NAME)"

test: ## Run tests (if any)
	@echo "Running tests..."
	go test -v ./...

docker-build: ## Build Docker image
	@echo "Building Docker image: $(IMAGE_NAME):$(IMAGE_TAG)..."
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .
	@echo "Docker image built successfully!"

docker-run: ## Run Docker container locally
	@echo "Running Docker container on port $(PORT)..."
	docker run --rm -p $(PORT):$(PORT) --name $(APP_NAME) $(IMAGE_NAME):$(IMAGE_TAG)

docker-stop: ## Stop running Docker container
	@echo "Stopping Docker container..."
	docker stop $(APP_NAME) || true

minikube-setup: ## Setup minikube environment
	@echo "Setting up minikube environment..."
	@echo "Starting minikube (if not running)..."
	minikube start || true
	@echo "Configuring Docker to use minikube's daemon..."
	@eval $$(minikube docker-env)

minikube-build: ## Build Docker image in minikube
	@echo "Building Docker image in minikube..."
	@eval $$(minikube docker-env) && docker build -t $(IMAGE_NAME):$(IMAGE_TAG) .
	@echo "Docker image built in minikube successfully!"

k8s-deploy: ## Deploy to Kubernetes (minikube)
	@echo "Deploying to Kubernetes..."
	kubectl apply -f k8s/deployment.yaml
	kubectl apply -f k8s/service.yaml
	@echo "Deployment complete!"
	@echo "Waiting for pods to be ready..."
	kubectl wait --for=condition=ready pod -l app=$(APP_NAME) --timeout=60s || true
	@echo ""
	@echo "To access the service, run: make k8s-url"

k8s-delete: ## Delete Kubernetes resources
	@echo "Deleting Kubernetes resources..."
	kubectl delete -f k8s/service.yaml || true
	kubectl delete -f k8s/deployment.yaml || true
	@echo "Resources deleted!"

k8s-status: ## Check Kubernetes deployment status
	@echo "=== Deployments ==="
	kubectl get deployments -l app=$(APP_NAME)
	@echo ""
	@echo "=== Pods ==="
	kubectl get pods -l app=$(APP_NAME)
	@echo ""
	@echo "=== Services ==="
	kubectl get services -l app=$(APP_NAME)
	@echo ""
	@echo "=== Events ==="
	kubectl get events --sort-by=.metadata.creationTimestamp | tail -10

k8s-logs: ## Show logs from Kubernetes pods
	@echo "Fetching logs from pods..."
	kubectl logs -l app=$(APP_NAME) --tail=50 -f

k8s-url: ## Get the service URL
	@echo "Service URL:"
	@minikube service $(APP_NAME) --url

k8s-open: ## Open service in browser
	@echo "Opening service in browser..."
	minikube service $(APP_NAME)

k8s-restart: ## Restart Kubernetes deployment
	@echo "Restarting deployment..."
	kubectl rollout restart deployment/$(APP_NAME)
	kubectl rollout status deployment/$(APP_NAME)

clean: ## Clean up build artifacts
	@echo "Cleaning up..."
	rm -f $(APP_NAME)
	@echo "Clean complete!"

full-deploy: minikube-build k8s-deploy ## Full deployment: build in minikube and deploy
	@echo ""
	@echo "=========================================="
	@echo "Full deployment complete!"
	@echo "=========================================="
	@echo ""
	@echo "Run 'make k8s-url' to get the service URL"
	@echo "Run 'make k8s-status' to check deployment status"
	@echo "Run 'make k8s-logs' to view logs"

quick-test: ## Quick test of all endpoints locally
	@echo "Testing endpoints..."
	@echo ""
	@echo "1. Testing root endpoint..."
	@curl -s http://localhost:$(PORT)/ | jq . || echo "Failed"
	@echo ""
	@echo "2. Testing health endpoint..."
	@curl -s http://localhost:$(PORT)/health | jq . || echo "Failed"
	@echo ""
	@echo "3. Testing info endpoint..."
	@curl -s http://localhost:$(PORT)/api/info | jq . || echo "Failed"
	@echo ""
	@echo "4. Testing echo endpoint..."
	@curl -s "http://localhost:$(PORT)/api/echo?message=Hello" | jq . || echo "Failed"
	@echo ""
	@echo "5. Testing data endpoint..."
	@curl -s -X POST http://localhost:$(PORT)/api/data \
		-H "Content-Type: application/json" \
		-d '{"name":"test","value":"data"}' | jq . || echo "Failed"