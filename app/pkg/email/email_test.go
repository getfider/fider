package email_test

import (
	"testing"

	"github.com/getfider/fider/app/pkg/log/noop"
	"github.com/getfider/fider/app/pkg/worker"

	"github.com/getfider/fider/app/pkg/email"

	. "github.com/getfider/fider/app/pkg/assert"
)

func TestRenderMessage(t *testing.T) {
	RegisterT(t)

	ctx := worker.NewContext("ID-1", "TaskName", nil, noop.NewLogger())
	message := email.RenderMessage(ctx, "echo_test", email.Params{
		"name": "Fider",
	})
	Expect(message.Subject).Equals("Message to: Fider")
	Expect(message.Body).Equals(`<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
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
			
			<tr>
				<td></td>
				<td>
					<table width="100%" bgcolor="#FFFFFF" cellpadding="0" cellspacing="0" border="0" style="text-align:left;padding:20px;border-radius:5px;color:#1c262d;border:1px solid #ECECEC;">
						Hello World Fider!
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
</html>`)
}

func TestEmailWhitelist_Valid(t *testing.T) {
	RegisterT(t)

	email.SetWhitelist("(^.+@fider.io$)|(^darthvader\\.fider(\\+.*)?@gmail\\.com$)")

	for _, address := range []string{
		"me@fider.io",
		"john@fider.io",
		"me+123@fider.io",
		"darthvader.fider@gmail.com",
		"darthvader.fider+1234@gmail.com",
		"darthvader.fider+434@gmail.com",
	} {
		Expect(email.CanSendTo(address)).IsTrue()
	}

	for _, address := range []string{
		"",
		"me@fidero.io",
		"@fider.io",
		"me+123@fider.iod",
		"darthvader@gmail.com",
		"darthvader.fider@gmaila.com",
		"darthvader.fidera@gmail.com",
	} {
		Expect(email.CanSendTo(address)).IsFalse()
	}
}

func TestParamsMerge(t *testing.T) {
	RegisterT(t)

	p1 := email.Params{
		"name": "Jon",
		"age":  26,
	}
	p2 := p1.Merge(email.Params{
		"age":   30,
		"email": "john.snow@got.com",
	})
	Expect(p2).Equals(email.Params{
		"name":  "Jon",
		"age":   30,
		"email": "john.snow@got.com",
	})
}
