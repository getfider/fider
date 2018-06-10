subject: {{ .title }}
body:

{{ if .duplicate }}
  <i>This idea has been closed as a <strong>{{ .status }}</strong> of {{ .duplicate }}.</i>
{{ else }}
  <i>Status has changed to <strong>{{ .status }}</strong>.</i> <br />

  {{ .content }}
{{ end }}

<span style="color:#666;font-size:11px">
— <br />
You are receiving this because you are subscribed to this thread. Please do not reply to this email. <br />
{{ .view }}, {{ .unsubscribe }} or {{ .change }}.
</span>