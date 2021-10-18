package cmd

import "github.com/getfider/fider/app/models/dto"

type SendMail struct {
	From         dto.Recipient
	To           []dto.Recipient
	TemplateName string
	Props        dto.Props
}
