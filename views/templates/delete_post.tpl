subject: [{{ .tenantName }}] {{ .title }}
body:
<tr>
  <td>
    <p style="padding-bottom:10px;border-bottom:1px solid #efefef;color:#1c262d">
      <strong>{{ .title }}</strong> has been <strong>deleted</strong>.
    </p>
    {{ .content }}
    <p style="color:#666;font-size:14px">
      — <br />
      You are receiving this email because you are subscribed to this post. You can {{ .change }}.
    </p>
  </td>
</tr>