package generator

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/ihkokil/launchpad/cli/pkg/spec"
)

const helmReleaseTpl = `apiVersion: helm.toolkit.fluxcd.io/v2
kind: HelmRelease
metadata:
  name: {{ .Name }}
  namespace: default
  labels:
    platform.launchpad.io/managed-by: launchpad-cli
    platform.launchpad.io/owner: {{ .Owner }}
    platform.launchpad.io/tier: {{ .Tier }}
spec:
  interval: 5m
  chart:
    spec:
      chart: ./charts/app
      sourceRef:
        kind: GitRepository
        name: flux-system
        namespace: flux-system
  install:
    remediation:
      retries: 3
  upgrade:
    remediation:
      retries: 3
  values:
    owner: {{ .Owner }}
    tier: {{ .Tier }}
    replicaCount: {{ if eq .Tier "critical" }}3{{ else }}2{{ end }}
    image:
      repository: {{ if .Image.Repository }}{{ .Image.Repository }}{{ else }}"{{ .Name }}"{{ end }}
      tag: {{ if .Image.Tag }}"{{ .Image.Tag }}"{{ else }}"latest-stable"{{ end }}
    service:
      port: {{ .Port }}
      targetPort: {{ .Port }}
    resources:
      requests:
        cpu: "{{ .Resources.CPU }}"
        memory: "{{ .Resources.Memory }}"
      limits:
        cpu: "500m"
        memory: "512Mi"
    {{- if .Env }}
    env:
      {{- range $k, $v := .Env }}
      {{ $k }}: "{{ $v }}"
      {{- end }}
    {{- end }}
    autoscaling:
      enabled: true
      minReplicas: {{ if eq .Tier "critical" }}3{{ else }}2{{ end }}
      maxReplicas: {{ if eq .Tier "critical" }}10{{ else }}5{{ end }}
    pdb:
      enabled: true
      minAvailable: 1
`

const slothSLOTpl = `apiVersion: sloth.slok.dev/v1
kind: PrometheusServiceLevel
metadata:
  name: {{ .Name }}-slo
  namespace: default
  labels:
    platform.launchpad.io/owner: {{ .Owner }}
    platform.launchpad.io/tier: {{ .Tier }}
spec:
  service: {{ .Name }}
  labels:
    owner: {{ .Owner }}
    tier: {{ .Tier }}
  slos:
    - name: availability
      objective: {{ .SLO.Availability }}
      description: "99.9% of HTTP requests respond with non-5xx status codes over a 30-day rolling window."
      sli:
        events:
          errorQuery: sum(rate(http_requests_total{job="{{ .Name }}",status=~"5.."}[{{"{{"}}.window{{"}}"}}]))
          totalQuery: sum(rate(http_requests_total{job="{{ .Name }}"}[{{"{{"}}.window{{"}}"}}]))
      alerting:
        name: {{ .Name }}-HighErrorRate
        labels:
          severity: page
          team: {{ .Owner }}
        annotations:
          summary: "High error budget burn rate on {{ .Name }}"
          runbook_url: "https://runbooks.launchpad.io/services/{{ .Name }}/availability"
          dashboard_url: "http://localhost:3000/d/launchpad-slo/{{ .Name }}"
        pageAlert:
          labels:
            severity: page
        ticketAlert:
          labels:
            severity: warning
    - name: latency-p95
      objective: 99.0
      description: "95th percentile latency must remain below {{ .SLO.LatencyP95MS }}ms."
      sli:
        events:
          errorQuery: sum(rate(http_request_duration_seconds_bucket{job="{{ .Name }}",le=~"0.[0-9]+"}[{{"{{"}}.window{{"}}"}}])) - sum(rate(http_request_duration_seconds_bucket{job="{{ .Name }}",le="{{ printf "%.2f" (div (float64 .SLO.LatencyP95MS) 1000.0) }}"}[{{"{{"}}.window{{"}}"}}]))
          totalQuery: sum(rate(http_request_duration_seconds_count{job="{{ .Name }}"}[{{"{{"}}.window{{"}}"}}]))
      alerting:
        name: {{ .Name }}-HighLatency
        labels:
          severity: warning
          team: {{ .Owner }}
        annotations:
          summary: "p95 latency exceeded threshold of {{ .SLO.LatencyP95MS }}ms on {{ .Name }}"
          runbook_url: "https://runbooks.launchpad.io/services/{{ .Name }}/latency"
`

func div(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

func RenderHelmRelease(s *spec.ServiceSpec) ([]byte, error) {
	tmpl, err := template.New("helmrelease").Parse(helmReleaseTpl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HelmRelease template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, s); err != nil {
		return nil, fmt.Errorf("failed to execute HelmRelease template: %w", err)
	}
	return buf.Bytes(), nil
}

func RenderSlothSLO(s *spec.ServiceSpec) ([]byte, error) {
	funcMap := template.FuncMap{
		"div":     div,
		"float64": func(i int) float64 { return float64(i) },
	}
	tmpl, err := template.New("sloth").Funcs(funcMap).Parse(slothSLOTpl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Sloth template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, s); err != nil {
		return nil, fmt.Errorf("failed to execute Sloth template: %w", err)
	}
	return buf.Bytes(), nil
}

