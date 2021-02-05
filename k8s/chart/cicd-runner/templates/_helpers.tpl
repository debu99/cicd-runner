{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "cicd-runner.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "cicd-runner.fullname" -}}
{{-   if .Values.fullnameOverride -}}
{{-     .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{-   else -}}
{{-     $name := default .Chart.Name .Values.nameOverride -}}
{{-     printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{-   end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "cicd-runner.chart" -}}
{{-   printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Define the name of the secret containing the tokens
*/}}
{{- define "cicd-runner.secret" -}}
{{- default (include "cicd-runner.fullname" .) .Values.runners.secret | quote -}}
{{- end -}}

{{/*
Define the name of the s3 cache secret
*/}}
{{- define "cicd-runner.cache.secret" -}}
{{- if .Values.runners.cache.secretName -}}
{{- .Values.runners.cache.secretName | quote -}}
{{- end -}}
{{- end -}}

{{/*
Template for outputing the gitlabUrl
*/}}
{{- define "cicd-runner.gitlabUrl" -}}
{{- .Values.gitlabUrl | quote -}}
{{- end -}}

{{/*
Template runners.cache.s3ServerAddress in order to allow overrides from external charts.
*/}}
{{- define "cicd-runner.cache.s3ServerAddress" }}
{{- default "" .Values.runners.cache.s3ServerAddress | quote -}}
{{- end -}}

{{/*
Define the image, using .Chart.AppVersion and CICD Runner image as a default value
*/}}
{{- define "cicd-runner.image" }}
{{-   $appVersion := ternary "bleeding" (print "v" .Chart.AppVersion) (eq .Chart.AppVersion "bleeding") -}}
{{-   $image := printf "debu99/cicd-runner.alpine-%s" $appVersion -}}
{{-   default $image .Values.image }}
{{- end -}}

{{/*
Unregister runners on pod stop
*/}}
{{- define "cicd-runner.unregisterRunners" -}}
{{- if or (and (hasKey .Values "unregisterRunners") .Values.unregisterRunners) (and (not (hasKey .Values "unregisterRunners")) .Values.runnerRegistrationToken) -}}
lifecycle:
  preStop:
    exec:
      command: ["/entrypoint", "unregister", "--all-runners"]
{{- end -}}
{{- end -}}
