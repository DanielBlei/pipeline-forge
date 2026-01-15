# Pipeline Forge 🔥 — Jenkins CI/CD Lab Branch
⚠️ **Experimental / Lab Branch**

This branch (`lab/jenkins-k8s`) exists to evaluate **Jenkins-based CI/CD** for Pipeline-Forge,
running **Jenkins on Kubernetes** with:

- Jenkins Helm chart
- Jenkins Configuration as Code (JCasC)
- Job DSL seed patterns
- Kubernetes-native agents

This work is intentionally **isolated from `main`** and is not part of the default CI
(GitHub Actions remains the primary CI for the project).

👉 **Start here:** [jenkins/README.md](./jenkins/README.md)

---

## What this lab demonstrates

- Running Jenkins locally on Kubernetes (kind / minikube)
- Bootstrapping pipelines via JCasC + Job DSL
- Separating infrastructure bootstrap from pipeline logic
- Keeping CI logic repo-driven and reproducible
