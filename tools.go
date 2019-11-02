// +build tools

package tools

import (
	_ "github.com/magefile/mage"
	_ "github.com/joho/godotenv/cmd/godotenv"
	_ "github.com/cosmtrek/air"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)