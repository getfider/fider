package oauth

const (
	//FacebookProvider is const for 'facebook'
	FacebookProvider = "facebook"
	//GoogleProvider is const for 'google'
	GoogleProvider = "google"
)

//UserProfile represents an OAuth user profile
type UserProfile struct {
	ID    string
	Name  string
	Email string
}

//Service provides OAuth operations
type Service interface {
	GetAuthURL(provider string, redirect string) string
	GetProfile(provider string, code string) (*UserProfile, error)
}
