package models

// AppSettings is an application-wide settings
type AppSettings struct {
	Mode            string `json:"mode"`
	BuildTime       string `json:"buildTime"`
	Version         string `json:"version"`
	Environment     string `json:"environment"`
	GoogleAnalytics string `json:"googleAnalytics"`
	Compiler        string `json:"compiler"`
}
