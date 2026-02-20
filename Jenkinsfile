pipeline {
    agent any

    parameters {
        file(name: 'CONFIG_FILE', description: 'Configuration file to pass to the binary')
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
                    go build -tags netgo -ldflags '-extldflags "-static"' ./cmd/gosha -o ./${BINARY_NAME}
                '''
            }
        }

        stage('Run Binary') {
            steps {
                echo 'Running binary with config file...'
                sh '''
                    mkdir -p input
                    cp "${CONFIG_FILE}" input/config.yaml
                    ./${BINARY_NAME} input/config.yaml
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