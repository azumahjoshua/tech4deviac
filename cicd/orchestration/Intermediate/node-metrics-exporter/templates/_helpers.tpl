{{/* Generate basic labels */}}
{{- define "node-metrics-exporter.labels" -}}
app.kubernetes.io/name: {{ include "node-metrics-exporter.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/* Selector labels */}}
{{- define "node-metrics-exporter.selectorLabels" -}}
app.kubernetes.io/name: {{ include "node-metrics-exporter.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/* Chart name */}}
{{- define "node-metrics-exporter.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end }}

{{/* Full name */}}
{{- define "node-metrics-exporter.fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end }}