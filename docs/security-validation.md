# Launchpad: Security Validation & Kyverno Guardrails Proof

**Author:** Md. Iqbal Haider Khan  
**Scope:** Automated Admission Control & Zero-Trust Enforcement  
**Cluster:** Multi-node `kind` (`v1.31.0`)  
**Engine:** Kyverno `v3.2.0+` & Kubernetes Dynamic Admission Webhooks  

---

## 1. Executive Summary

In a modern Internal Developer Platform (IDP), security cannot rely on manual code reviews or post-deployment scans alone. Launchpad implements **shift-left dynamic admission control** using Kyverno to automatically reject non-compliant workloads at the Kubernetes API server boundary before a single container can be scheduled.

---

## 2. Active Security Policies

All policies run in `validationFailureAction: Enforce` mode across all developer workloads:

| Policy Name | Severity | Enforced Rule | Rationale |
| :--- | :--- | :--- | :--- |
| **`disallow-root-user`** | High | `runAsNonRoot: true` & `runAsUser > 0` | Blocks container breakout and host filesystem tampering. |
| **`require-resource-limits`** | Medium | CPU & Memory `requests` and `limits` | Prevents noisy-neighbor starvation and node OOM crashes. |
| **`disallow-latest-tag`** | Medium | Blocks `:latest` or untagged images | Enforces immutable git commit SHAs or semver tags for deterministic rollbacks. |
| **`zero-trust-networkpolicy`** | High | Default-deny ingress/egress | Restricts pod ingress strictly to ingress controller and monitoring stack. |

---

## 3. Negative Testing Proof Artifact (Workload Rejected)

To prove that guardrails are actively enforced, an intentionally insecure manifest ([`test/bad-workload.yaml`](file:///home/ihkokil/CV/launchpad/test/bad-workload.yaml)) was applied to the cluster:

```yaml
# test/bad-workload.yaml
apiVersion: v1
kind: Pod
metadata:
  name: bad-insecure-workload
  namespace: default
spec:
  securityContext:
    runAsNonRoot: false
    runAsUser: 0
  containers:
    - name: insecure-container
      image: nginx:latest
      # Omitting CPU/memory requests & limits
```

### Execution Log (API Server Rejection)

```bash
$ kubectl apply -f test/bad-workload.yaml
```

**Captured Admission Webhook Response:**

```text
Error from server: error when creating "test/bad-workload.yaml": admission webhook "validate.kyverno.svc-fail" denied the request: 

resource Pod/default/bad-insecure-workload was blocked due to the following policies 

disallow-latest-tag:
  validate-image-tag: 'validation error: Using '':latest'' or untagged images is prohibited. Use a semver or commit SHA tag. rule validate-image-tag failed at path /spec/containers/0/image/'
disallow-root-user:
  check-runAsNonRoot: 'validation error: Running as root is forbidden. Set ''securityContext.runAsNonRoot: true'' and ''runAsUser'' > 0. rule check-runAsNonRoot failed at path /spec/securityContext/runAsNonRoot/'
require-resource-limits:
  check-resources: 'validation error: CPU and Memory requests and limits are strictly required on all containers. rule check-resources failed at path /spec/containers/0/resources/limits/'
```

---

## 4. Positive Testing Proof Artifact (Golden Path Compliant)

A microservice deployed through the Launchpad Golden Path Helm Chart (`charts/app`) automatically inherits compliant security contexts, immutable tags, probes, and resource constraints:

```bash
$ helm template compliant-app charts/app \
    --set image.repository="compliant-app" \
    --set image.tag="v1.0.0" | kubectl apply -f -
```

**Output:**
```text
networkpolicy.networking.k8s.io/compliant-app-launchpad-app created
poddisruptionbudget.policy/compliant-app-launchpad-app created
serviceaccount/compliant-app-launchpad-app created
service/compliant-app-launchpad-app created
deployment.apps/compliant-app-launchpad-app created
horizontalpodautoscaler.autoscaling/compliant-app-launchpad-app created
servicemonitor.monitoring.coreos.com/compliant-app-launchpad-app created
```

**Result:** Zero warnings, 100% compliance out-of-the-box. Developers are protected by default without having to manually research Kubernetes security benchmarks.

