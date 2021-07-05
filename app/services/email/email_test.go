package email_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/services/email"

	. "github.com/getfider/fider/app/pkg/assert"
)

func TestRenderMessage(t *testing.T) {
	RegisterT(t)

	message := email.RenderMessage(context.Background(), "echo_test", dto.Props{
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
	<body bgcolor="#F7F7F7" style="font-size:18px">
		<table width="100%" bgcolor="#F7F7F7" cellpadding="0" cellspacing="0" border="0" style="text-align:center;font-size:18px;">
			<tr>
				<td height="40">&nbsp;</td>
			</tr>
			
			<tr>
				<td align="center">
					<table bgcolor="#FFFFFF" cellpadding="0" cellspacing="0" border="0" style="text-align:left;padding:20px;margin:10px;border-radius:5px;color:#1c262d;border:1px solid #ECECEC;min-width:320px;max-width:660px;">
						
Hello World Fider!

					</table>
				</td>
			</tr>
			<tr>
				<td>
					<span style="color:#666;font-size:12px">This email was sent from a notification-only address that cannot accept incoming email. Please do not reply to this message.</span>
				</td>
			</tr>
			<tr>
				<td height="40">&nbsp;</td>
			</tr>
		</table>
	</body>
</html>`)
}

func TestCanSendTo(t *testing.T) {
	RegisterT(t)

	testCases := []struct {
		allowlist string
		blocklist string
		input     []string
		canSend   bool
	}{
		{
			allowlist: "(^.+@fider.io$)|(^darthvader\\.fider(\\+.*)?@gmail\\.com$)",
			blocklist: "",
			input:     []string{"me@fider.io", "me+123@fider.io", "darthvader.fider@gmail.com", "darthvader.fider+434@gmail.com"},
			canSend:   true,
		},
		{
			allowlist: "(^.+@fider.io$)|(^darthvader\\.fider(\\+.*)?@gmail\\.com$)",
			blocklist: "",
			input:     []string{"me+123@fider.iod", "me@fidero.io", "darthvader.fidera@gmail.com", "@fider.io"},
			canSend:   false,
		},
		{
			allowlist: "(^.+@fider.io$)|(^darthvader\\.fider(\\+.*)?@gmail\\.com$)",
			blocklist: "(^.+@fider.io$)",
			input:     []string{"me@fider.io"},
			canSend:   true,
		},
		{
			allowlist: "",
			blocklist: "(^.+@fider.io$)",
			input:     []string{"me@fider.io", "abc@fider.io"},
			canSend:   false,
		},
		{
			allowlist: "",
			blocklist: "(^.+@fider.io$)",
			input:     []string{"me@fider.com", "abc@fiderio.io"},
			canSend:   true,
		},
		{
			allowlist: "",
			blocklist: "",
			input:     []string{"me@fider.io"},
			canSend:   true,
		},
		{
			allowlist: "",
			blocklist: "",
			input:     []string{"", " "},
			canSend:   false,
		},
	}

	for _, testCase := range testCases {
		email.SetAllowlist(testCase.allowlist)
		email.SetBlocklist(testCase.blocklist)
		for _, input := range testCase.input {
			Expect(email.CanSendTo(input)).Equals(testCase.canSend)
		}
	}
}

func TestRecipient_String(t *testing.T) {
	RegisterT(t)

	testCases := []struct {
		name     string
		email    string
		expected string
	}{
		{
			name:     "Jon",
			email:    "jon@got.com",
			expected: `"Jon" <jon@got.com>`,
		},
		{
			name:     "Snow, Jon",
			email:    "jon@got.com",
			expected: `"Snow, Jon" <jon@got.com>`,
		},
		{
			name:     "",
			email:    "jon@got.com",
			expected: "<jon@got.com>",
		},
		{
			name:     "Jon's Home Account",
			email:    "jon@got.com",
			expected: `"Jon's Home Account" <jon@got.com>`,
		},
		{
			name:     `Jon "Great" Snow`,
			email:    "jon@got.com",
			expected: `"Jon \"Great\" Snow" <jon@got.com>`,
		},
		{
			name:     "Jon",
			email:    "",
			expected: "",
		},
	}

	for _, testCase := range testCases {
		r := dto.NewRecipient(testCase.name, testCase.email, dto.Props{})
		Expect(r.String()).Equals(testCase.expected)
	}
}
