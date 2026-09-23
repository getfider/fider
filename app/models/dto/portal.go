package dto

// PortalSummary is one entry in the public portal directory. URL and LogoURL are absolute and
// point at the portal's own host, because the directory is served from the root domain where
// tenant-scoped paths do not resolve.
type PortalSummary struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Host    string `json:"host"`
	LogoURL string `json:"logoURL,omitempty"`
}
