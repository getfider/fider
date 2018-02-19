subject: {{ .title }}
body:

{{ if .duplicate }}
  <p>This idea has been closed as a <strong>{{ .status }}</strong> of {{ .duplicate }}. </p>
{{ else }}
  Status has changed to <strong>{{ .status }}</strong>. <br />

  {{ .content }}
{{ end }}

<span style="color:#666;font-size:11px">
— <br />
You are receiving this because you are subscribed to this thread. Please do not reply to this email. <br />
{{ .view }}, {{ .unsubscribe }} or {{ .change }}.
</span>