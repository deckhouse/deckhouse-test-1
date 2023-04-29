{{- define "webhook.labels" }}
{{- include "labels" (list . (dict "webhook" .Chart.Name)) }}
{{- end }}
