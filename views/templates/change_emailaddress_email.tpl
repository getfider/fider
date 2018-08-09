subject: Confirm your new email
body:
<tr>
  <td>
    <p>Hi <strong>{{ .name }}</strong>,</p>
    <p>You have requested to change your email from {{ .oldEmail }} to {{ .newEmail }}.</p>
    <p>Click the link below to confirm this operation.</p>
  </td>
</tr>
<tr>
  <td>{{ .link }}</td>
</tr>
<tr>
  <td>
    <span style="color:#666;font-size:11px">This link will expire in 24 hours and can only be used once.</span>
  </td>
</tr>