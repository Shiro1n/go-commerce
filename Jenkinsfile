pipeline {
    agent any

    environment {
        DOCKER_USERNAME = 'shiroin'
        DOCKER_AUTH_SERVICE_IMAGE = "${DOCKER_USERNAME}/auth-service:latest"
        DOCKER_USER_SERVICE_IMAGE = "${DOCKER_USERNAME}/user-service:latest"
    }

    stages {
        stage('Checkout') {
            steps {
                script {
                    echo "Checking out code from Git..."
                    checkout scm
                }
            }
        }
        stage('Build Docker Images') {
            steps {
                script {
                    echo "Building Docker images..."
                    sh 'docker build -t $DOCKER_AUTH_SERVICE_IMAGE ./auth-service'
                    sh 'docker build -t $DOCKER_USER_SERVICE_IMAGE ./user-service'
                }
            }
        }
        stage('Push Docker Images') {
            steps {
                script {
                    echo "Pushing Docker images to Docker Hub..."
                    sh 'docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD'
                    sh 'docker push $DOCKER_AUTH_SERVICE_IMAGE'
                    sh 'docker push $DOCKER_USER_SERVICE_IMAGE'
                }
            }
        }
        stage('Deploy to Kubernetes') {
            steps {
                script {
                    echo "Applying Kubernetes configurations..."

                    // Apply Redis configuration
                    sh 'kubectl apply -f ./k8s/redis/redis-pvc.yaml'
                    sh 'kubectl apply -f ./k8s/redis/redis-deployment.yaml'
                    sh 'kubectl apply -f ./k8s/redis/redis-service.yaml'

                    // Apply PostgreSQL configuration
                    sh 'kubectl apply -f ./k8s/postgres/postgres-pvc.yaml'
                    sh 'kubectl apply -f ./k8s/postgres/postgres-deployment.yaml'
                    sh 'kubectl apply -f ./k8s/postgres/postgres-service.yaml'

                    // Apply User Service configuration
                    sh 'kubectl apply -f ./user-service/k8s/user-service-deployment.yaml'
                    sh 'kubectl apply -f ./user-service/k8s/user-service-service.yaml'

                    // Apply Auth Service configuration
                    sh 'kubectl apply -f ./auth-service/k8s/auth-deployment.yaml'
                    sh 'kubectl apply -f ./auth-service/k8s/auth-service.yaml'

                    // Apply cert-manager
                    sh 'kubectl apply -f ./gateway/cert-manager.yaml'

                    // Apply ClusterIssuer for Let's Encrypt
                    sh 'kubectl apply -f ./gateway/cluster-issuer.yaml'

                    // Apply Ingress configuration
                    sh 'kubectl apply -f ./gateway/ingress.yaml'
                }
            }
        }
    }
    post {
        success {
            script {
                echo 'Pipeline completed successfully!'
            }
        }
        failure {
            script {
                echo 'Pipeline failed.'
            }
        }
    }
}
