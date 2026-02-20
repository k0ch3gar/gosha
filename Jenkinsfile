pipeline {
    agent any

    parameters {
        stashedFile 'large'
    }

    tools {
        go 'goshaaa'
    }

    environment {
        GO111MODULE = 'on'
        CGO_ENABLED = 0
        BINARY_NAME = 'myapp'
    }

    stages {
        stage('Checkout') {
            steps {
                echo 'Cloning repository...'
                checkout scm
            }
        }

        stage('Build') {
            steps {
                echo 'Building Go binary...'
                sh '''
                    go mod tidy
                    go mod download
                    go build -o ./${BINARY_NAME} -tags netgo -ldflags '-extldflags "-static"' ./cmd/gosha
                '''
            }
        }

        stage('Run Binary') {
            steps {
                echo 'Running binary with example file...'
                unstash 'large'
                sh '''
                    ./${BINARY_NAME} large
                '''
            }
        }
    }

    post {
        always {
            echo 'Cleaning up workspace...'
            cleanWs()
            archiveArtifacts artifacts: 'bin/**', allowEmptyArchive: true
        }
        success {
            echo 'Pipeline succeeded!'
        }
        failure {
            echo 'Pipeline failed!'
        }
    }
}