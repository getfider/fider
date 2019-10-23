package cmd

import "github.com/getfider/fider/app/models/dto"

type SendMail struct {
	From         string
	To           []dto.Recipient
	TemplateName string
	Props        dto.Props
}
