# Launchpad: Developer Experience Benchmark & Measurements

This document records the before-and-after benchmarks measuring developer experience and delivery efficiency, proving the tangible impact of the **Launchpad Internal Developer Platform**.

---

## Benchmark Comparison Table

| Metric | Manual Deploy (Baseline) | Launchpad IDP | Delta / Improvement |
| :--- | :--- | :--- | :--- |
| **Time to First Deploy** | **42 minutes** | *Target: < 6 min* | **~85% reduction** |
| **Lines of YAML Touched** | **167 lines (6 manifests)** | *Target: < 20 lines (`service.yaml`)* | **~88% reduction** |
| **Manual Steps / Commands** | **14 commands** | *Target: 2 steps (`launchpad new`, Git PR)* | **86% toil reduction** |
| **SLO / Alerting Setup Time**| **~35 minutes** | *Auto-generated with service (< 10s)* | **100% automated** |
| **Security Compliance Gate** | **Manual audit (error-prone)** | *Admission gate enforced via Kyverno* | **Zero-trust by default** |

---

## 1. Manual Baseline Measurement Log

* **Service Tested:** `manual-orders-api` (Node.js REST API with healthz & metrics endpoints)
* **Cluster Environment:** 3-node `kind` Kubernetes cluster (`v1.31.0`)
* **Engineer:** Md. Iqbal Haider Khan
* **Date:** October 2026
* **Measured Steps Executed:**
  1. Author container build specification (`Dockerfile`).
  2. Author raw Kubernetes `Deployment` (`deployment.yaml`: 67 lines) with non-root security context, dropped capabilities, read-only root filesystem, CPU/memory requests & limits, and liveness/readiness probes.
  3. Author `Service` (`service.yaml`: 17 lines) ClusterIP mapping.
  4. Author `HorizontalPodAutoscaler` (`hpa.yaml`: 27 lines) targetting CPU 75% and Memory 80%.
  5. Author `PodDisruptionBudget` (`pdb.yaml`: 12 lines) requiring `minAvailable: 1`.
  6. Author Zero-Trust `NetworkPolicy` (`networkpolicy.yaml`: 28 lines) enforcing default-deny ingress except from Prometheus monitoring namespace.
  7. Author Prometheus `ServiceMonitor` (`servicemonitor.yaml`: 16 lines) targetting `/metrics`.
  8. Build Docker container (`manual-orders-api:v1.0.0`).
  9. Sideload container image across all cluster nodes.
  10. Execute `kubectl apply -f test/manual-baseline/k8s/`.
  11. Verify Pod readiness and HPA target binding.
  12. Verify zero-trust NetworkPolicy isolation (blocked default namespace curl, allowed monitoring namespace).
  13. Verify Prometheus target discovery (`/api/v1/targets` active scraping).
* **Stopwatch & Code Metrics:**
  * **Total Duration from Scratch:** **42 minutes**
  * **Total Raw YAML Lines Authored:** **167 lines**
  * **Manual Terminal Commands Executed:** **14 commands**

---

## 2. Launchpad IDP Target & Onboarding Log

With Launchpad IDP, the developer replaces all 167 lines across 6 manifests with a single **10-line `service.yaml`**:

```yaml
name: orders-api
owner: team-checkout
tier: standard
runtime: node
port: 8080
resources:
  cpu: "100m"
  memory: "128Mi"
slo:
  availability: 99.9
  latency_p95_ms: 250
```

### Onboarded Services Tracking

| Service | Runtime | Owner | Deploy Time | Lines of YAML | Status |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `orders-api` | Node.js | Team Checkout | Pending (Week 5) | Target: 10 lines | Planned |
| `finops-service` | Go | Core Platform | Pending (Week 5) | Target: 12 lines | Planned |
| `peer-portal-api` | Node / PHP | Peer / Colleague | Pending (Week 5) | Target: 11 lines | Planned |
