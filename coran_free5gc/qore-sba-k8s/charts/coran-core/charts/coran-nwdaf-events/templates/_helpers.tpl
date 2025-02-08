

{{/*
Expand the name of the chart.
*/}}
{{- define "coran-nwdaf.name" -}}
{{- default .Chart.Name .Values.nwdaf.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "coran-nwdaf.fullname" -}}
{{- if .Values.nwdaf.fullnameOverride }}
{{- .Values.nwdaf.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nwdaf.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "coran-nwdaf.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "coran-nwdaf.labels" -}}
helm.sh/chart: {{ include "coran-nwdaf.chart" . }}
{{ include "coran-nwdaf.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "coran-nwdaf.selectorLabels" -}}
app.kubernetes.io/name: {{ include "coran-nwdaf.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
NWDaf Pod Annotations
*/}}
{{- define "coran-nwdaf.nwdafAnnotations" -}}
{{- with .Values.nwdaf }}
{{- if .podAnnotations }}
{{- toYaml .podAnnotations }}
{{- end }}
{{- end }}
{{- end }}
