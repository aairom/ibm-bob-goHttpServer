# Go HTTP Server - Architecture Diagrams

This document contains comprehensive Mermaid diagrams illustrating the architecture and workflows of the Go HTTP Server project.

## 1. System Architecture

```mermaid
graph TB
    subgraph "Client Layer"
        Client[Client/Browser]
    end
    
    subgraph "Kubernetes Cluster (Minikube)"
        Service[Service<br/>NodePort :30080]
        
        subgraph "Deployment (2 Replicas)"
            Pod1[Pod 1<br/>go-http-server]
            Pod2[Pod 2<br/>go-http-server]
        end
    end
    
    subgraph "Container"
        App[Go HTTP Server<br/>:8080]
        Health[Health Checks<br/>Liveness/Readiness]
    end
    
    Client -->|HTTP Request| Service
    Service -->|Load Balance| Pod1
    Service -->|Load Balance| Pod2
    Pod1 --> App
    Pod2 --> App
    App -.->|Probe| Health
    
    style Client fill:#e1f5ff
    style Service fill:#fff4e1
    style Pod1 fill:#e8f5e9
    style Pod2 fill:#e8f5e9
    style App fill:#f3e5f5
    style Health fill:#fce4ec
```

## 2. API Endpoints Flow

```mermaid
graph LR
    Client[Client] --> Router{HTTP Router}
    
    Router -->|GET /| Home[Home Handler<br/>Welcome Message]
    Router -->|GET /health| Health[Health Handler<br/>Status & Uptime]
    Router -->|GET /api/info| Info[Info Handler<br/>Server Info]
    Router -->|GET /api/echo| Echo[Echo Handler<br/>Message Echo]
    Router -->|POST /api/data| Data[Data Handler<br/>JSON Processing]
    
    Home --> MW1[CORS Middleware]
    Health --> MW2[CORS Middleware]
    Info --> MW3[CORS Middleware]
    Echo --> MW4[CORS Middleware]
    Data --> MW5[CORS Middleware]
    
    MW1 --> Log1[Logging Middleware]
    MW2 --> Log2[Logging Middleware]
    MW3 --> Log3[Logging Middleware]
    MW4 --> Log4[Logging Middleware]
    MW5 --> Log5[Logging Middleware]
    
    Log1 --> Resp1[JSON Response]
    Log2 --> Resp2[JSON Response]
    Log3 --> Resp3[JSON Response]
    Log4 --> Resp4[JSON Response]
    Log5 --> Resp5[JSON Response]
    
    style Router fill:#fff4e1
    style Home fill:#e8f5e9
    style Health fill:#e8f5e9
    style Info fill:#e8f5e9
    style Echo fill:#e8f5e9
    style Data fill:#e8f5e9
```

## 3. Deployment Workflow

```mermaid
sequenceDiagram
    participant Dev as Developer
    participant Docker as Docker
    participant Minikube as Minikube
    participant K8s as Kubernetes
    participant Pod as Pod
    
    Dev->>Docker: Build Image<br/>(multi-stage)
    Docker->>Docker: Stage 1: Build Go Binary
    Docker->>Docker: Stage 2: Create Alpine Image
    Docker-->>Dev: Image Ready
    
    Dev->>Minikube: eval $(minikube docker-env)
    Dev->>Docker: Build in Minikube
    Docker-->>Minikube: Image Available
    
    Dev->>K8s: kubectl apply deployment.yaml
    K8s->>Pod: Create 2 Replicas
    Pod->>Pod: Run Health Checks
    Pod-->>K8s: Ready
    
    Dev->>K8s: kubectl apply service.yaml
    K8s->>K8s: Create NodePort Service
    K8s-->>Dev: Service Available :30080
    
    Dev->>Minikube: minikube service go-http-server
    Minikube-->>Dev: Open Browser
```

## 4. Request Processing Flow

```mermaid
flowchart TD
    Start([HTTP Request]) --> CORS{CORS<br/>Middleware}
    CORS -->|Add Headers| Log{Logging<br/>Middleware}
    Log -->|Log Request| Route{Route<br/>Handler}
    
    Route -->|/| Home[Home Handler]
    Route -->|/health| Health[Health Handler]
    Route -->|/api/info| Info[Info Handler]
    Route -->|/api/echo| Echo[Echo Handler]
    Route -->|/api/data| Data[Data Handler]
    
    Home --> JSON1[Encode JSON]
    Health --> JSON2[Encode JSON]
    Info --> JSON3[Encode JSON]
    Echo --> Valid{Valid<br/>Query?}
    Data --> Method{POST<br/>Method?}
    
    Valid -->|Yes| JSON4[Encode JSON]
    Valid -->|No| Err1[400 Error]
    
    Method -->|Yes| Parse{Valid<br/>JSON?}
    Method -->|No| Err2[405 Error]
    
    Parse -->|Yes| JSON5[201 Created]
    Parse -->|No| Err3[400 Error]
    
    JSON1 --> End([Response])
    JSON2 --> End
    JSON3 --> End
    JSON4 --> End
    JSON5 --> End
    Err1 --> End
    Err2 --> End
    Err3 --> End
    
    style Start fill:#e1f5ff
    style End fill:#e1f5ff
    style CORS fill:#fff4e1
    style Log fill:#fff4e1
    style Route fill:#f3e5f5
```

## 5. Docker Multi-Stage Build

```mermaid
graph TB
    subgraph "Stage 1: Builder (golang:1.21-alpine)"
        S1[Copy go.mod] --> S2[Download Dependencies]
        S2 --> S3[Copy main.go]
        S3 --> S4[Build Static Binary<br/>CGO_ENABLED=0]
        S4 --> Binary[server binary<br/>~15MB]
    end
    
    subgraph "Stage 2: Runtime (alpine:latest)"
        R1[Install ca-certificates] --> R2[Create Non-Root User]
        R2 --> R3[Copy Binary from Builder]
        R3 --> R4[Set Permissions]
        R4 --> Final[Final Image<br/>~15MB]
    end
    
    Binary -.->|Copy| R3
    
    style Binary fill:#e8f5e9
    style Final fill:#e8f5e9
```

## 6. Application Lifecycle

```mermaid
stateDiagram-v2
    [*] --> Starting: Server Start
    Starting --> Initializing: Load Config
    Initializing --> RegisteringRoutes: Setup Routes
    RegisteringRoutes --> ApplyingMiddleware: Add Middleware
    ApplyingMiddleware --> Running: Listen on Port 8080
    
    Running --> Processing: Incoming Request
    Processing --> Running: Response Sent
    
    Running --> ShuttingDown: SIGINT/SIGTERM
    ShuttingDown --> GracefulShutdown: Wait for Requests (30s)
    GracefulShutdown --> [*]: Server Stopped
    
    Processing --> Error: Request Failed
    Error --> Running: Log Error
```

## Diagram Descriptions

### 1. System Architecture
Shows the complete deployment structure with Kubernetes components including the Service (NodePort), Deployment with 2 replicas, and the containerized Go application with health checks.

### 2. API Endpoints Flow
Illustrates how HTTP requests are routed to different handlers and processed through CORS and logging middleware layers.

### 3. Deployment Workflow
Sequence diagram showing the step-by-step process of building the Docker image and deploying to Kubernetes via Minikube.

### 4. Request Processing Flow
Detailed flowchart of how each HTTP request is processed, including validation, error handling, and response generation.

### 5. Docker Multi-Stage Build
Demonstrates the two-stage Docker build process that creates a minimal production image (~15MB).

### 6. Application Lifecycle
State diagram showing the server's lifecycle from startup through request processing to graceful shutdown.

---

**Generated for:** Go HTTP Server v1.0.0  
**Date:** 2025-12-14  
**Author:** Bob
