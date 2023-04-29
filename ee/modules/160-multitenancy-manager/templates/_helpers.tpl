{{- $moduleName := "multitenancy-manager" }}

{{- define "name" }}
{{- printf "name: %s-%s" (index . 0).Chart.Name (index . 1) }}
{{- end }}

{{- define "namespace" }}
{{- printf "namespace: d8-%s" .Chart.Name }}
{{- end }}

{{- define "name.namespace" }}
{{- include "name" . }}
{{ include "namespace" (index . 0) }}
{{- end }}

{{- define "labels" }}
{{- $context := index . 0 }}
{{- $additionalLabels := dict }}
{{- if gt ( len . ) 1 }}
{{- $additionalLabels = index . 1 }}
{{- end }}
{{- include "helm_lib_module_labels" (list $context ( merge (dict "app" "multitenancy-manager") $additionalLabels )) }}
{{- end }}
