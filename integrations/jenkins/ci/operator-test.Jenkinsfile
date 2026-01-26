pipeline {
  agent {
    kubernetes {
      label 'golang-agent'
      defaultContainer 'golang'
      yaml """
apiVersion: v1
kind: Pod
metadata:
  labels:
    app: jenkins-golang
spec:
  serviceAccountName: jenkins
  restartPolicy: Never
  containers:
    - name: golang
      image: golang:1.24-bookworm
      command: ["cat"]
      tty: true
"""
    }
  }

  options {
    timeout(time: 30, unit: 'MINUTES')
    timestamps()
    disableConcurrentBuilds()
    buildDiscarder(logRotator(numToKeepStr: '10'))
  }

  environment {
    WORKDIR = 'operator'
    GOLANGCI_LINT_VERSION = 'v2.5.0'
  }

  stages {
    stage('Checkout') {
      steps {
        checkout scmGit(
          branches: scm.branches,
          extensions: [cloneOption(shallow: true, depth: 1)],
          userRemoteConfigs: scm.userRemoteConfigs
        )
      }
    }

    stage('Run Linter') {
      steps {
        container('golang') {
          dir("${env.WORKDIR}") {
            sh '''
              set -eu
              export GOBIN="$PWD/bin"
              go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}
              "$GOBIN"/golangci-lint version
              "$GOBIN"/golangci-lint run --config ../.golangci.yml
            '''
          }
        }
      }
    }

    stage('Verify Generated Code') {
      steps {
        container('golang') {
          dir("${env.WORKDIR}") {
            sh '''
              set -eu
              export GOBIN="$PWD/bin"
              go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.18.0
              "$GOBIN"/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
              if ! git diff --exit-code; then
                echo "Generated code is not up to date. Please run 'make generate' and commit changes."
                exit 1
              fi
            '''
          }
        }
      }
    }

    stage('Go Mod Tidy Check') {
      steps {
        container('golang') {
          dir("${env.WORKDIR}") {
            sh '''
              set -eu
              go mod tidy
              if ! git diff --exit-code; then
                echo "go.mod is not tidy. Please run 'go mod tidy' and commit changes."
                exit 1
              fi
            '''
          }
        }
      }
    }

    stage('Run Tests') {
      steps {
        container('golang') {
          dir("${env.WORKDIR}") {
            sh '''
              set -eu
              export GOBIN="$PWD/bin"

              # Basic checks similar to Makefile deps
              go fmt ./...
              go vet ./...

              # Run unit tests (exclude e2e)
              go test $(go list ./... | grep -v /e2e) -coverprofile cover.out
            '''
          }
        }
      }
      post {
        always {
          archiveArtifacts artifacts: "${env.WORKDIR}/cover.out", allowEmptyArchive: true
        }
      }
    }
  }

  post {
    always {
      echo 'Operator test pipeline completed'
    }
    success {
      echo 'Build succeeded!'
    }
    failure {
      echo 'Build failed!'
    }
  }
}
