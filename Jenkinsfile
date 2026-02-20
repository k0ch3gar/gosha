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
                    go build -o ./${BINARY_NAME} -tags netgo -ldflags '-extldflags "-static"' ./cmd/gosha
                '''
            }
        }

        stage('Run Binary') {
            steps {
                echo 'Running binary with example file...'
                script {
                    def configFile = params.CONFIG_FILE
                    if (configFile) {
                        writeFile file: 'input/input.sh', text: readFile(configFile)
                        echo "Example file written to input/input.sh"
                    } else {
                        error "CONFIG_FILE parameter is required!"
                   }
                }
                sh '''
                    ./${BINARY_NAME} input/input.sh
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