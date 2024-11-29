pipeline {
    agent {
        node {
            label 'agent'
        }
    }
    environment {
        DOCKER_IMAGE = 'crypticseeds/secure-api-management-platform'
    }
    tools {
        go 'Go'
        dockerTool 'docker'
    }
    stages {
        stage ('Git Checkout') {
            steps {
                git branch: 'main', credentialsId: 'git-cred', url: 'https://github.com/crypticseeds/secure-api-management-platform.git'
            }
        }
        stage ('Install Dependencies') {
            steps {
                sh 'go mod tidy || exit 1'
            }
        }
        stage ('Run Tests') {
            steps {
                sh '''
                    go test ./... -v > test-results.txt || exit 1
                '''
            }
        }
        stage ('Trivy fs Scan') {
            steps {
                sh '''
                    trivy fs --format table \
                    --exit-code 1 \
                    --severity HIGH,CRITICAL \
                    -o fs-report.html .
                '''
            }
        }
        stage ('Sonaqube Scan') {
            steps {
                script {
                    def scannerHome = tool 'sonar-scanner' 
                    withSonarQubeEnv('sonar') {
                        sh "${scannerHome}/bin/sonar-scanner -Dsonar.projectKey=secure-api-management-platform -Dsonar.projectName=secure-api-management-platform -Dsonar.sources=."
                    }
                }
            }
        }
        stage ('Build Docker Image & Tag') {
            steps {
                script {
                    withDockerRegistry(credentialsId: 'docker-cred') {
                        sh """
                            docker build --no-cache \
                            -t ${env.DOCKER_IMAGE}:latest \
                            -t ${env.DOCKER_IMAGE}:v1.${BUILD_NUMBER} .
                        """
                    }
                }
            }
        }
        stage ('Trivy Image Scan') {
            steps {
                sh """
                    trivy image --format table \
                    --exit-code 1 \
                    --severity HIGH,CRITICAL \
                    --scanners vuln \
                    -o image-report.html ${env.DOCKER_IMAGE}:latest
                """
            }
        }
        stage ('Push Docker Image') {
            steps {
                script {
                    withDockerRegistry(credentialsId: 'docker-cred') {
                        sh """
                            docker push ${env.DOCKER_IMAGE}:latest
                            docker push ${env.DOCKER_IMAGE}:v1.${BUILD_NUMBER}
                        """
                    }
                }
            }
        }
    }
    post {
        always {
            cleanWs()
            sh 'docker system prune -f'
        }
    }
}
