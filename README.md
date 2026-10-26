# Launchpad 🚀
### Production-Grade Internal Developer Platform (IDP)

[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://go.dev/)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-v1.31-326CE5.svg)](https://kubernetes.io/)
[![Flux v2](https://img.shields.io/badge/GitOps-Flux_v2-29B6F6.svg)](https://fluxcd.io/)
[![Kyverno](https://img.shields.io/badge/Policy-Kyverno-00B2B2.svg)](https://kyverno.io/)
[![Terraform](https://img.shields.io/badge/IaC-Terraform-7B42BC.svg)](https://www.terraform.io/)

Launchpad is an Internal Developer Platform designed to eliminate developer friction and enforce zero-trust security and reliability defaults. Instead of writing hundreds of lines of Kubernetes YAML, cloud infrastructure, and Prometheus alerting rules, product engineers interact with a single 10-line contract (`service.yaml`).

---

## 🏗️ Architecture

```
+----------------------------------------------------------------------------------------------------+
| DEVELOPER EXPERIENCE                                                                               |
|                                                                                                    |
|  1. Developer runs:                                                                                |
|     $ launchpad new orders-api --runtime node                                                      |
|        │                                                                                           |
|        ▼ (Go CLI scaffolds repo from template + opens PR in platform repo with service.yaml)       |
+--------┼-------------------------------------------------------------------------------------------+
         │
         ▼
+----------------------------------------------------------------------------------------------------+
| PLATFORM AUTOMATION (GitHub Actions & Terraform)                                                    |
|                                                                                                    |
|  2. Platform CI:                                                                                   |
|     - Validates service.yaml against JSON schema                                                   |
|     - Runs `launchpad render` to generate K8s manifests & HelmRelease                              |
|     - Runs `terraform plan` (ECR repo + IAM OIDC role with least privilege)                        |
|  3. Service CI:                                                                                    |
|     - Runs tests -> Trivy container vulnerability scan -> Multi-stage non-root build               |
|     - Pushes image to ECR -> Opens PR to GitOps repo updating image tag                            |
+--------┼-------------------------------------------------------------------------------------------+
         │
         ▼
+----------------------------------------------------------------------------------------------------+
| CLOUD PLATFORM & KUBERNETES RUNTIME (AWS EKS & Flux v2)                                            |
|                                                                                                    |
|  4. Flux GitOps Controller reconciles desired state from GitOps repo                               |
|  5. Kyverno Admission Controller rejects non-compliant workloads:                                  |
|     - Block root user / privilege escalation                                                       |
|     - Enforce CPU/Memory requests & limits                                                         |
|     - Block `:latest` image tag                                                                    |
|  6. "Golden Path" Helm Chart applies platform defaults:                                            |
|     - Probes (liveness/readiness/startup), HPA, PDB, TopologySpreadConstraints                     |
|     - Default-Deny NetworkPolicies (only allow explicit ingress & egress)                          |
+--------┼-------------------------------------------------------------------------------------------+
         │
         ▼
+----------------------------------------------------------------------------------------------------+
| OBSERVABILITY & RELIABILITY (Sloth, Prometheus, Grafana & ChatOps)                                 |
|                                                                                                    |
|  7. Ingress/Gateway emits RED metrics (Rate, Errors, Duration)                                     |
|  8. Sloth generates multi-window multi-burn-rate alerts (99.9% availability, p95 < 250ms)            |
|  9. Alerts fire on symptoms (error budget depletion), routing to Slack/Teams with runbook links    |
+--------┼-------------------------------------------------------------------------------------------+
         │
         ▼
+----------------------------------------------------------------------------------------------------+
| CHAOS GAME DAY & RESILIENCE (Chaos Mesh + k6)                                                      |
|                                                                                                    |
|  10. Continuous synthetic traffic (k6) + Injected failures (Chaos Mesh):                           |
|      - Pod kills, node drains, dependency latency, network partition                               |
|  11. Published Game Day Postmortem: Time-to-detect, time-to-recover, error budget burned           |
+----------------------------------------------------------------------------------------------------+
```

---

## 📑 Developer Contract: `service.yaml`

```yaml
name: orders-api
owner: team-checkout
tier: standard          # critical (3+ replicas, strict PDB) | standard (2 replicas)
runtime: node           # node | go | php
port: 8080
env:
  NODE_ENV: production
resources:
  cpu: "100m"
  memory: "128Mi"
slo:
  availability: 99.9    # % success over 30 days
  latency_p95_ms: 250   # 95th percentile under 250ms
```

---

## 📂 Repository Structure

```text
launchpad/
├── cli/                            # Go CLI tool (`launchpad new`, `render`, `status`)
├── infra/                          # Terraform code (AWS VPC, EKS, ECR, IAM OIDC)
├── charts/app/                     # Reusable "Golden Path" Helm chart
├── policies/kyverno/               # Security guardrails & NetworkPolicies
├── gitops/
│   └── clusters/
│       └── local/                  # Local kind cluster Flux manifests
├── services/                       # Developer manifests (orders-api.yaml, etc.)
├── chaos/                          # Chaos Mesh tests & Game Day postmortem reports
└── docs/
    └── measurements.md             # Baseline vs Platform benchmark numbers
```

---

## 👤 Author & Maintainer
**Md. Iqbal Haider Khan**  
B.Sc. in Computer Science and Engineering, RUET  
Platform & DevOps Engineer
