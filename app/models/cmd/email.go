package cmd

import "github.com/Spicy-Bush/fider-tarkov-community/app/models/dto"

type SendMail struct {
	From         dto.Recipient
	To           []dto.Recipient
	TemplateName string
	Props        dto.Props
}
