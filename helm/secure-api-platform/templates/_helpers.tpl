{{/*
Expand the name of the chart.
*/}}
{{- define "secure-api-platform.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "secure-api-platform.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "secure-api-platform.labels" -}}
helm.sh/chart: {{ include "secure-api-platform.chart" . }}
{{ include "secure-api-platform.selectorLabels" . }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "secure-api-platform.selectorLabels" -}}
app.kubernetes.io/name: {{ include "secure-api-platform.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Database environment variables
*/}}
{{- define "secure-api-platform.databaseEnv" -}}
- name: POSTGRES_DB
  value: {{ .Values.postgresql.database }}
- name: POSTGRES_USER
  valueFrom:
    secretKeyRef:
      name: {{ include "secure-api-platform.fullname" . }}-secrets
      key: POSTGRES_USER
- name: POSTGRES_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ include "secure-api-platform.fullname" . }}-secrets
      key: POSTGRES_PASSWORD
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "secure-api-platform.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}