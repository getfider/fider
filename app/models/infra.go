package models

// SystemSettings is the system-wide settings
type SystemSettings struct {
	Mode            string `json:"mode"`
	BuildTime       string `json:"buildTime"`
	Version         string `json:"version"`
	Environment     string `json:"environment"`
	GoogleAnalytics string `json:"googleAnalytics"`
	Compiler        string `json:"compiler"`
	Domain          string `json:"domain"`
}
