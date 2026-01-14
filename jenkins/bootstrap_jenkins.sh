#!/usr/bin/env bash

set -euo pipefail

K8S_RUNTIME=${1:-kind}
ROLLOUT_TIMEOUT=5m
NAMESPACE=jenkins
RELEASE_CHART=jenkins

function apply_k8s_resources() {
  echo "Creating Kubernetes resources (Service Account, RBAC, PV, PVC)..."
  kubectl create ns "$NAMESPACE"
  kubectl apply -n "$NAMESPACE" -f k8s/jenkins-serviceAccount-and-rbac.yaml
  kubectl apply -n "$NAMESPACE" -f k8s/jenkins-volume.yaml
}

function helm_install() {
    echo "Installing Jenkins via Helm..."
    helm repo add jenkinsci https://charts.jenkins.io > /dev/null
    helm repo update >/dev/null
    helm install "$RELEASE_CHART" -n "$NAMESPACE" -f values.lab.yaml jenkinsci/jenkins
}

function  wait_rollout() {
  echo "Waiting for rollout to complete..."
  if ! kubectl rollout status statefulset/jenkins -n "$NAMESPACE" --timeout=$ROLLOUT_TIMEOUT; then
    echo "Rollout timed out, double check the pod/container logs and events for more details." >&2
    exit 1
  fi
}

function main() {
  if [[ $K8S_RUNTIME == "kind" ]]; then
    kind create cluster --name jenkins # suffix -control-plane is added by kind
  elif [[ $K8S_RUNTIME == "minikube" ]]; then
    minikube start -p jenkins-control-plane
  else
    echo "Unknown runtime: $K8S_RUNTIME"
    exit 1
  fi

  apply_k8s_resources
  helm_install
  wait_rollout

  printf '\033[1;32m%s\033[0m\n\n' '✔ Jenkins is ready!'
  # shellcheck disable=SC1083
  echo -e "admin password: $(kubectl get secret -n "jenkins" jenkins -o jsonpath={.data.jenkins-admin-password} | base64 --decode)\n"
  echo "Run 'kubectl port-forward svc/jenkins -n \"$NAMESPACE\" 8080:8080' and access Jenkins UI on localhost:8080"
}

main "$@"