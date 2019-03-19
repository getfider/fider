package cmd

import "github.com/getfider/fider/app/models/dto"

type LogDebug struct {
	Message string
	Props   dto.Props
}

type LogError struct {
	Err     error
	Message string
	Props   dto.Props
}

type LogWarn struct {
	Message string
	Props   dto.Props
}

type LogInfo struct {
	Message string
	Props   dto.Props
}
