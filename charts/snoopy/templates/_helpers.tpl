{{/*
Generate snoopyApiToken if not specified in values
*/}}
{{ define "common.snoopy-api-token" }}
{{- if .Values.common.snoopyApiToken }}
  {{- .Values.common.snoopyApiToken -}}
{{- else if (lookup "v1" "Secret" .Release.Namespace "common-snoopy-secret").data }}
  {{- $obj := (lookup "v1" "Secret" .Release.Namespace "common-snoopy-secret").data -}}
  {{- index $obj "SNOOPY_API_TOKEN" | b64dec -}}
{{- else -}}
  {{- randAlphaNum 48 -}}
{{- end -}}
{{- end -}}
{{/*
Generate ingress hostname if not specified in values
*/}}
{{- define "ingress.hostName" -}}
{{- .Values.snoopy.ingress.domain.prefix }}{{ if ne .Values.snoopy.ingress.domain.prefix "" }}.{{ end }}{{- .Values.snoopy.ingress.domain.base }}{{- .Values.snoopy.ingress.domain.suffix }}
{{- end }}
