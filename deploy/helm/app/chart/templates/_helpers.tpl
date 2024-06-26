{{- define "envvars" -}}
{{- range $key, $value := .Values.environmentVariables }}
- name: {{ $key }}
  value: {{ $value | quote }}
{{- end }}
{{- end }}
