package models

// AppSettings is an application-wide settings
type AppSettings struct {
	BuildTime       string
	Version         string
	AuthEndpoint    string
	Environment     string
	GoogleAnalytics string
	Compiler        string
}
