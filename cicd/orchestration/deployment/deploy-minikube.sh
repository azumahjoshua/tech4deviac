#!/bin/bash

# Start Minikube if not running
if ! minikube status | grep -q "Running"; then
  echo "Starting Minikube..."
  minikube start
fi

# Function to deploy a component
deploy() {
  local deployment_file="$1"
  local service_file="$2"

  echo "Deploying $deployment_file and $service_file..."
  kubectl apply -f "$deployment_file"
  kubectl apply -f "$service_file"
}

# Deploy banking-api
deploy "banking-api-deployment.yaml" "banking-api-service.yaml"

# Deploy banking-processor
deploy "banking-processor-deployment.yaml" "banking-processor-service.yaml"

# Deploy banking-frontend
deploy "banking-frontend-deployment.yaml" "banking-frontend-service.yaml"

# List all services
kubectl get svc

# If LoadBalancer is used, expose them via NodePort (Minikube workaround)
echo "Checking for LoadBalancer services..."
for svc in banking-api banking-processor banking-frontend; do
  if kubectl get svc "$svc" -o jsonpath='{.spec.type}' | grep -q "LoadBalancer"; then
    echo "Exposing $svc with NodePort..."
    kubectl patch svc "$svc" -p '{"spec": {"type": "NodePort"}}'
  fi
done

echo "Deployment completed successfully!"
