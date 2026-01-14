
# Jenkins CI/CD Lab for Pipeline-Forge

This directory contains a local Jenkins setup for experimenting with CI/CD Jenkins pipelines for the Pipeline-Forge project.

## Overview

This lab environment provides a quick way to spin up Jenkins locally using Kubernetes (kind or minikube) and configure it with Jenkins Configuration as Code (JCasc).

The goal is to test and develop CI/CD pipelines for Pipeline-Forge components, evaluating Jenkins solutions and features.

## Quick Start

### Prerequisites

- `kubectl`
- `helm`
- `docker`
- Either `kind` or `minikube`

### Installation

Run the bootstrap script to create a local Jenkins instance:

```bash
# Kind by default
./bootstrap_jenkins.sh

# or specify minikube
./bootstrap_jenkins.sh minikube
```

The script will:
1. Create a local Kubernetes cluster
2. Set up the `jenkins` namespace with required RBAC and storage
3. Install Jenkins via Helm with JCasc configuration
4. Wait for Jenkins to be ready


After the installation completes:

1. **Port forward to Jenkins:**
   ```bash
   kubectl port-forward svc/jenkins -n jenkins 8080:8080
   ```

2. **Get the admin password:**
   ```bash
   kubectl get secret -n jenkins jenkins -o jsonpath={.data.jenkins-admin-password} | base64 --decode
   ```

## Bootstrapping Pipelines

**⚠️ Important:** After Jenkins is installed, you need to manually run the seed job to bootstrap all pipelines.

**Steps:**
1. Navigate to Jenkins UI
2. Locate and run the seed job (configured via JCasc)
3. This will create all Pipeline-Forge CI/CD pipelines (work in progress)


## Cleanup

Delete the local Jenkins cluster:
```bash
kind delete cluster --name jenkins
# or
minikube delete -p jenkins-control-plane
```

## Lab Implementation Notes

Compared to the [upstream Jenkins Helm chart defaults](https://raw.githubusercontent.com/jenkinsci/helm-charts/main/charts/jenkins/values.yaml), this lab includes:

**Seed Job via JCasC**
A seed pipeline job is defined via Jenkins Configuration as Code
- The seed job pulls the Pipeline-Forge repository and bootstraps Jenkins CI pipelines from repository-managed definitions

**Two-Step Bootstrap (Intentional Design)**
- Jenkins installation and CI/CD pipeline creation are deliberately separated
- Keeps bootstrap logs clean and avoids hiding configuration issues
- Allows you to debug seed job errors without re-running the entire cluster setup
- Enables independent iteration on CI/CD pipeline definitions

**CI-Focused Plugin Set**
- Additional plugins for multibranch pipelines, GitHub integration, Job DSL, Kubernetes agents, and JCasC

**Job DSL Script Security**
- Script security for Job DSL is disabled via init script to allow seed job execution without needing manual approval
- Can be manually re-enabled in Jenkins UI (Manage Jenkins → Configure System → Job DSL) after the initial seed job completes if required

**Local-First Defaults**
- Jenkins URL: `http://localhost:8080`
- Access via port-forwarding with no external service

**Configuration**

Override values are defined in `values.lab.yaml`.

**Note:** Plugin versions are intentionally unpinned for this development environment to always pull the latest versions.

### References

- [Installing Jenkins on Kubernetes](https://www.jenkins.io/doc/book/installing/kubernetes/#install-jenkins)
- [Jenkins Helm Chart Default Values](https://raw.githubusercontent.com/jenkinsci/helm-charts/main/charts/jenkins/values.yaml)