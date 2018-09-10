{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "dex.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "dex.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "dex.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "dex.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "dex.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Create the health check path
*/}}

{{/*
Create secret key from environment variables
*/}}
{{- define "dex.envkey" -}}
{{ . | replace "_" "-" | lower }}
{{- end -}}

{{/*
Create the healthCheckPath for readiness and liveness probes.

Based on the following template values:
    - healthCheckPath
    - ingress.path

The default is '/healthz'
*/}}

{{- define "dex.healthCheckPath" -}}
{{- if .Values.healthCheckPath -}}
  {{ .Values.healthCheckPath }}
{{- else -}}
  {{- if .Values.ingress.enabled -}}
    {{ default "" .Values.ingress.path | trimSuffix "/" }}/healthz
  {{- else -}}
    {{ default "/healthz" }}
  {{- end -}}
{{- end -}}
{{- end -}}
