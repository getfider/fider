<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <meta name="viewport" content="width=device-width">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
  </head>
  <body bgcolor="#F7F7F7" style="font-size:16px">
    <table width="100%" bgcolor="#F7F7F7" cellpadding="0" cellspacing="0" border="0" style="text-align:center;font-size:14px;">
      <tr>
        <td width="200"></td>
        <td height="40">&nbsp;</td>
        <td width="200"></td>
      </tr>
      {{ if .logo }}
        <tr>
          <td></td>
          <td>
            <img height="50" src="{{ .logo }}"/>
          </td>
          <td></td>
        </tr>
        <tr>
          <td></td>
          <td height="10" style="line-height:1px;">&nbsp;</td>
          <td></td>
        </tr>
      {{ end }}
      <tr>
        <td></td>
        <td>
          <table width="100%" bgcolor="#FFFFFF" cellpadding="0" cellspacing="0" border="0" style="text-align:left;padding:20px;border-radius:5px;color:#1c262d;border:1px solid #ECECEC;">
            {{ .body }}
          </table>
        </td>
        <td></td>
      </tr>
      <tr>
        <td></td>
        <td height="10" style="line-height:1px;">&nbsp;</td>
        <td></td>
      </tr>
      <tr>
        <td></td>
        <td>
          <span style="color:#666;font-size:11px">This email was sent from a notification-only address that cannot accept incoming email. Please do not reply to this message.</span>
        </td>
        <td></td>
      </tr>
      <tr>
        <td></td>
        <td height="40">&nbsp;</td>
        <td></td>
      </tr>
    </table>
  </body>
</html>