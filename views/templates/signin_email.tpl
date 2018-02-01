subject: Sign in to {{ .tenantName }}
body:
Click the link below to sign in. 
<br /><br />
<a href='{{ .baseURL }}/signin/verify?k={{ .verificationKey }}'>{{ .baseURL }}/signin/verify?k={{ .verificationKey }}</a> 
<br /><br />
<span style="color:#b3b3b1;font-size:11px">This link will expire in 15 minutes and can only be used once.</span>