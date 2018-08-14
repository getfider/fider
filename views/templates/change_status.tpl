subject: [{{ .tenantName }}] {{ .title }}
body:
<tr>
  <td>
    {{ if .duplicate }}
      <strong>{{ .title }} ({{ .postLink }})</strong> has been closed as a <strong>{{ .status }}</strong> of {{ .duplicate }}.
    {{ else }}
      Status of <strong>{{ .title }} ({{ .postLink }})</strong> has changed to <strong>{{ .status }}</strong>.
    {{ end }}
  </td>
</tr>
<tr>
  <td></td>
  <td height="10" style="line-height:1px;">&nbsp;</td>
  <td></td>
</tr>
<tr>
  <td style="border-top:1px solid #efefef;">{{ .content }}</td>
</tr>
<tr>
  <td>
    <span style="color:#666;font-size:11px">
    — <br />
    You are receiving this because you are subscribed to this thread. <br />
    {{ .view }}, {{ .unsubscribe }} or {{ .change }}.
    </span>
  </td>
</tr>