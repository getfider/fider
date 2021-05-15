package mailgun_test

import (
	"context"
	"io/ioutil"
	"net/url"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/services/email/mailgun"
	"github.com/getfider/fider/app/services/httpclient/httpclientmock"

	"github.com/getfider/fider/app/services/email"

	. "github.com/getfider/fider/app/pkg/assert"
)

var ctx context.Context

func reset() {
	ctx = context.WithValue(context.Background(), app.TenantCtxKey, &entity.Tenant{
		Subdomain: "got",
	})
	bus.Init(mailgun.Service{}, httpclientmock.Service{})
}

func TestSend_Success(t *testing.T) {
	RegisterT(t)
	env.Config.HostMode = "multi"
	reset()

	bus.Publish(ctx, &cmd.SendMail{
		From: "Fider Test",
		To: []dto.Recipient{
			{
				Name:    "Jon Sow",
				Address: "jon.snow@got.com",
			},
		},
		TemplateName: "echo_test",
		Props: dto.Props{
			"name": "Hello",
		},
	})

	Expect(httpclientmock.RequestsHistory).HasLen(1)
	Expect(httpclientmock.RequestsHistory[0].URL.String()).Equals("https://api.mailgun.net/v3/mydomain.com/messages")
	Expect(httpclientmock.RequestsHistory[0].Header.Get("Authorization")).Equals("Basic YXBpOm15czNjcjN0azN5")
	Expect(httpclientmock.RequestsHistory[0].Header.Get("Content-Type")).Equals("application/x-www-form-urlencoded")

	bytes, err := ioutil.ReadAll(httpclientmock.RequestsHistory[0].Body)
	Expect(err).IsNil()
	values, err := url.ParseQuery(string(bytes))
	Expect(err).IsNil()
	Expect(values).HasLen(6)
	Expect(values.Get("to")).Equals(`"Jon Sow" <jon.snow@got.com>`)
	Expect(values.Get("from")).Equals(`"Fider Test" <noreply@random.org>`)
	Expect(values.Get("h:Reply-To")).Equals("noreply@random.org")
	Expect(values.Get("subject")).Equals("Message to: Hello")
	Expect(values["o:tag"][0]).Equals("template:echo_test")
	Expect(values["o:tag"][1]).Equals("tenant:got")
	Expect(values.Get("html")).Equals(`<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
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
						Hello World Hello!
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

func TestSend_SkipEmptyAddress(t *testing.T) {
	RegisterT(t)
	reset()

	bus.Publish(ctx, &cmd.SendMail{
		From: "Fider Test",
		To: []dto.Recipient{
			{
				Name:    "Jon Sow",
				Address: "",
			},
		},
		TemplateName: "echo_test",
		Props: dto.Props{
			"name": "Hello",
		},
	})

	Expect(httpclientmock.RequestsHistory).HasLen(0)
}

func TestSend_SkipUnlistedAddress(t *testing.T) {
	RegisterT(t)
	reset()
	email.SetAllowlist("^.*@gmail.com$")

	bus.Publish(ctx, &cmd.SendMail{
		From: "Fider Test",
		To: []dto.Recipient{
			{
				Name:    "Jon Sow",
				Address: "jon.snow@got.com",
			},
		},
		TemplateName: "echo_test",
		Props: dto.Props{
			"name": "Hello",
		},
	})

	Expect(httpclientmock.RequestsHistory).HasLen(0)
}

func TestBatch_Success(t *testing.T) {
	RegisterT(t)
	reset()
	email.SetAllowlist("")

	bus.Publish(ctx, &cmd.SendMail{
		From: "Fider Test",
		To: []dto.Recipient{
			{
				Name:    "Jon Sow",
				Address: "jon.snow@got.com",
				Props: dto.Props{
					"name": "Jon",
				},
			},
			{
				Name:    "Arya Stark",
				Address: "arya.start@got.com",
				Props: dto.Props{
					"name": "Arya",
				},
			},
		},
		TemplateName: "echo_test",
	})

	Expect(httpclientmock.RequestsHistory).HasLen(1)
	Expect(httpclientmock.RequestsHistory[0].URL.String()).Equals("https://api.mailgun.net/v3/mydomain.com/messages")
	Expect(httpclientmock.RequestsHistory[0].Header.Get("Authorization")).Equals("Basic YXBpOm15czNjcjN0azN5")
	Expect(httpclientmock.RequestsHistory[0].Header.Get("Content-Type")).Equals("application/x-www-form-urlencoded")

	bytes, err := ioutil.ReadAll(httpclientmock.RequestsHistory[0].Body)
	Expect(err).IsNil()
	values, err := url.ParseQuery(string(bytes))
	Expect(err).IsNil()
	Expect(values).HasLen(7)
	Expect(values["to"]).HasLen(2)
	Expect(values["to"][0]).Equals(`"Jon Sow" <jon.snow@got.com>`)
	Expect(values["to"][1]).Equals(`"Arya Stark" <arya.start@got.com>`)
	Expect(values.Get("from")).Equals(`"Fider Test" <noreply@random.org>`)
	Expect(values.Get("h:Reply-To")).Equals("noreply@random.org")
	Expect(values.Get("subject")).Equals("Message to: %recipient.name%")
	Expect(values["o:tag"]).HasLen(2)
	Expect(values["o:tag"][0]).Equals("template:echo_test")
	Expect(values["o:tag"][1]).Equals("tenant:got")
	Expect(values.Get("recipient-variables")).Equals("{\"arya.start@got.com\":{\"name\":\"Arya\"},\"jon.snow@got.com\":{\"name\":\"Jon\"}}")
	Expect(values.Get("html")).Equals(`<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
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
						Hello World %recipient.name%!
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

func TestGetBaseURL(t *testing.T) {
	RegisterT(t)
	reset()

	sendMail := &cmd.SendMail{
		From: "Fider Test",
		To: []dto.Recipient{
			{
				Name:    "Jon Sow",
				Address: "jon.snow@got.com",
			},
		},
		TemplateName: "echo_test",
		Props: dto.Props{
			"name": "Hello",
		},
	}

	// Fall back to US if there is nothing set
	env.Config.Email.Mailgun.Region = ""
	bus.Publish(ctx, sendMail)
	Expect(httpclientmock.RequestsHistory[0].URL.String()).Equals("https://api.mailgun.net/v3/mydomain.com/messages")

	// Return the EU domain for EU, ignore the case
	env.Config.Email.Mailgun.Region = "EU"
	bus.Publish(ctx, sendMail)
	Expect(httpclientmock.RequestsHistory[1].URL.String()).Equals("https://api.eu.mailgun.net/v3/mydomain.com/messages")

	env.Config.Email.Mailgun.Region = "eu"
	bus.Publish(ctx, sendMail)
	Expect(httpclientmock.RequestsHistory[2].URL.String()).Equals("https://api.eu.mailgun.net/v3/mydomain.com/messages")

	// Return the US domain for US, ignore the case
	env.Config.Email.Mailgun.Region = "US"
	bus.Publish(ctx, sendMail)
	Expect(httpclientmock.RequestsHistory[3].URL.String()).Equals("https://api.mailgun.net/v3/mydomain.com/messages")
	env.Config.Email.Mailgun.Region = "us"
	bus.Publish(ctx, sendMail)
	Expect(httpclientmock.RequestsHistory[4].URL.String()).Equals("https://api.mailgun.net/v3/mydomain.com/messages")

	// Return the US domain if the region is invalid
	env.Config.Email.Mailgun.Region = "Mars"
	bus.Publish(ctx, sendMail)
	Expect(httpclientmock.RequestsHistory[5].URL.String()).Equals("https://api.mailgun.net/v3/mydomain.com/messages")

}
