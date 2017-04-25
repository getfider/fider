package models

// WeCHYSettings is an application-wide settings
type WeCHYSettings struct {
	BuildTime    string
	Version      string
	AuthEndpoint string
	Environment  string
	Compiler     string
}
