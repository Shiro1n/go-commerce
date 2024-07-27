#!/bin/bash

# Define variables
dockerUsername="shiroin"
authServiceImage="$dockerUsername/auth-service:latest"
userServiceImage="$dockerUsername/user-service:latest"

# Build Docker images
echo "Building Docker images..."
docker build -t $authServiceImage ./auth-service
docker build -t $userServiceImage ./user-service

# Push Docker images to Docker Hub
echo "Pushing Docker images to Docker Hub..."
docker push $authServiceImage
docker push $userServiceImage

# Apply Kubernetes configurations
echo "Applying Kubernetes configurations..."

# Apply Redis configuration
kubectl apply -f ./k8s/redis/redis-pvc.yaml
kubectl apply -f ./k8s/redis/redis-deployment.yaml
kubectl apply -f ./k8s/redis/redis-service.yaml

# Apply PostgreSQL configuration
kubectl apply -f ./k8s/postgres/postgres-pvc.yaml
kubectl apply -f ./k8s/postgres/postgres-deployment.yaml
kubectl apply -f ./k8s/postgres/postgres-service.yaml

# Apply User Service configuration
kubectl apply -f ./user-service/k8s/user-service-deployment.yaml
kubectl apply -f ./user-service/k8s/user-service-service.yaml

# Apply Auth Service configuration
kubectl apply -f ./auth-service/k8s/auth-deployment.yaml
kubectl apply -f ./auth-service/k8s/auth-service.yaml

echo "Deployment completed successfully!"
