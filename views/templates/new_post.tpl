subject: [{{ .tenantName }}] {{ .title }}
body:
<tr>
  <td>
    <strong>{{ .userName }}</strong> created a new post <strong>{{ .title }} ({{ .postLink }})</strong>.
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
    You are receiving this because you are subscribed to this event. <br />
    {{ .view }} or {{ .change }}.
    </span>
  </td>
</tr>