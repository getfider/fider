subject: Confirm your new e-mail
body:
Hi {{ .name }},
<br /><br />
Looks like you have requested to change your e-mail from {{ .oldEmail }} to {{ .newEmail }}.
<br />
Click the link below to confirm this operation.
<br /><br />
<a href='{{ .baseURL }}/change-email/verify?k={{ .verificationKey }}'>{{ .baseURL }}/change-email/verify?k={{ .verificationKey }}</a> 
<br /><br />
<span style="color:#b3b3b1;font-size:11px">This link will expire in 24 hours and can only be used once.</span>