{{/*
Expand the name of the chart.
*/}}
{{- define "coran-nrf.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "coran-scp.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "coran-scp.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}{{- end }}
{{- end }}
{{- end }}

{{- define "coran-scp.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "coran-scp.labels" -}}
helm.sh/chart: {{ include "coran-scp.chart" . }}
{{ include "coran-scp.selectorLabels" . }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{- define "coran-scp.selectorLabels" -}}
app.kubernetes.io/name: {{ include "coran-scp.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "coran-scp.scpAnnotations" -}}
{{- with .Values.scp }}
{{- if .podAnnotations }}
{{- toYaml .podAnnotations | nindent 8 }}
{{- end }}
{{- end }}
{{- end }}

{{/*
NRF Pod Annotations
*/}}
{{- define "coran-nrf.nrfAnnotations" -}}
{{- with .Values.nrf }}
{{- if .podAnnotations }}
{{- toYaml .podAnnotations }}
{{- end }}
{{- end }}
{{- end }}
