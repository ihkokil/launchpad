{{/*
Expand the name of the app.
*/}}
{{- define "launchpad-app.name" -}}
{{- default .Release.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "launchpad-app.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- default .Release.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "launchpad-app.labels" -}}
helm.sh/chart: {{ include "launchpad-app.name" . }}-{{ .Chart.Version | replace "+" "_" }}
{{ include "launchpad-app.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
platform.launchpad.io/tier: {{ .Values.tier | default "standard" | quote }}
platform.launchpad.io/owner: {{ .Values.owner | default "platform-team" | quote }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "launchpad-app.selectorLabels" -}}
app.kubernetes.io/name: {{ include "launchpad-app.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

