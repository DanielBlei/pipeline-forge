pipeline {
  agent {
    kubernetes {
      label 'podman-agent'
      defaultContainer 'tools'
      yaml """
apiVersion: v1
kind: Pod
metadata:
  labels:
    app: jenkins-podman
spec:
  serviceAccountName: jenkins
  restartPolicy: Never
  containers:
    - name: tools
      image: quay.io/containers/podman:latest
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
    IMAGE_TAG = "quay.io/danielblei/pipeline-forge/ingest:ci-${env.GIT_COMMIT}"
    WORKDIR = 'workloads/ingest'
    STORAGE_DRIVER = 'vfs'
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

    stage('Setup Python Environment') {
      steps {
        dir("${env.WORKDIR}") {
          sh '''
            set -euo pipefail

            # Install UV package manager
            curl -LsSf https://astral.sh/uv/install.sh | sh
            export PATH="$HOME/.local/bin:$PATH"

            # Verify we're in the right directory with pyproject.toml
            ls -la pyproject.toml

            # Need review here to install right python version
            # Create virtual environment and install dependencies from pyproject.toml
            uv venv
            uv sync --frozen
            uv pip install -e . --no-deps
          '''
        }
      }
    }

    stage('Quality Checks') {
      parallel {
        stage('Ruff Linter') {
          steps {
            dir("${env.WORKDIR}") {
              sh '''
                set -euo pipefail
                export PATH="$HOME/.local/bin:$PATH"
                uv run ruff check .
              '''
            }
          }
        }

        stage('Ruff Format Check') {
          steps {
            dir("${env.WORKDIR}") {
              sh '''
                set -euo pipefail
                export PATH="$HOME/.local/bin:$PATH"
                uv run ruff format --check .
              '''
            }
          }
        }

        stage('Type Check (mypy)') {
          steps {
            dir("${env.WORKDIR}") {
              sh '''
                set -euo pipefail
                export PATH="$HOME/.local/bin:$PATH"
                uv run mypy .
              '''
            }
          }
        }
      }
    }

    stage('Run Tests') {
      steps {
        dir("${env.WORKDIR}") {
            sh '''
              set -euo pipefail
              export PATH="$HOME/.local/bin:$PATH"
              uv run pytest . -v --junitxml=test-results.xml
            '''
        }
      }
      post {
        always {
          // Publish test results if they exist
          junit allowEmptyResults: true, testResults: "${env.WORKDIR}/test-results.xml"
          archiveArtifacts artifacts: "${env.WORKDIR}/test-results.xml", allowEmptyArchive: true
        }
      }
    }

    stage('Build Image (Podman)') {
      steps {
        dir("${env.WORKDIR}") {
          container('tools') {
            sh '''
              set -euo pipefail
              podman --storage-driver="$STORAGE_DRIVER" build -t "$IMAGE_TAG" .
            '''
          }
        }
      }
    }

    stage('Smoke Test (uv run)') {
      steps {
        dir("${env.WORKDIR}") {
          sh '''
            set -euo pipefail
            export PATH="$HOME/.local/bin:$PATH"
            uv run ingest --help
          '''
        }
      }
    }
  }

  post {
    always {
      echo 'Pipeline completed'
    }
    success {
      echo 'Build succeeded!'
    }
    failure {
      echo 'Build failed!'
    }
    cleanup {
      container('tools') {
        sh '''
          podman --storage-driver="$STORAGE_DRIVER" rmi -f "$IMAGE_TAG" 2>/dev/null || true
        '''
      }
    }
  }
}
