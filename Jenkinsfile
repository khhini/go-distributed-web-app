pipeline {
    agent any
    tools {
        go 'go1.18'
    }
    environment {
        GO114MODULE = 'on'
        CGO_ENABLED = 0 
        GOPATH = "${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"
        GIN_MODE="release"
        DOCKERHUB_CREDENTIALS=credentials('dockerhub')
        GIT_COMMIT = sh(returnStdout: true, script: "git rev-parse --short=10 HEAD").trim()
    }
    stages {
        stage('Test Source Code'){
            steps {
                sh 'go mod tidy'
                sh 'go test'
            }
        }
        stage('Build Docker Image'){
            steps {
                sh 'docker build -t khhini/go-recipes-api:${GIT_COMMIT} .'
            }
        }
        stage('Push Docker Image'){
            steps {
                sh 'echo $DOCKERHUB_CREDENTIALS_PSW | docker login -u $DOCKERHUB_CREDENTIALS_USR --password-stdin'
                sh "docker push khhini/go-recipes-api:${GIT_COMMIT}"
                sh "docker image rm -f khhini/go-recipes-api:${GIT_COMMIT}"
                sh "docker system prune -f"
            }
        }
        stage('Deploy'){
            steps {
                ansiblePlaybook credentialsId: 'ansiblecd', disableHostKeyChecking: true, extras: ' -e IMAGE_TAG=${GIT_COMMIT}', inventory: 'ansible-deployment/hosts', playbook: 'ansible-deployment/pb_deployment.yml'
            }
        }
    }
}