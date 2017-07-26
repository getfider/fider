package oauth

//MockOAuthService implements a mocked OAuthService
type MockOAuthService struct{}

//GetAuthURL returns authentication url for given provider
func (p *MockOAuthService) GetAuthURL(authEndpoint string, provider string, redirect string) string {
	return "http://orange.test.canherayou.com/oauth/token?provider=" + provider + "&redirect=" + redirect
}

//GetProfile returns user profile based on provider and code
func (p *MockOAuthService) GetProfile(authEndpoint string, provider string, code string) (*UserProfile, error) {
	if provider == "facebook" && code == "123" {
		return &UserProfile{
			ID:    "FB123",
			Name:  "Jon Snow",
			Email: "jon.snow@got.com",
		}, nil
	}

	if provider == "facebook" && code == "456" {
		return &UserProfile{
			ID:    "FB456",
			Name:  "Some Facebook Guy",
			Email: "some.guy@facebook.com",
		}, nil
	}

	if provider == "facebook" && code == "798" {
		return &UserProfile{
			ID:    "FB798",
			Name:  "Mark",
			Email: "",
		}, nil
	}

	if provider == "google" && code == "123" {
		return &UserProfile{
			ID:    "GO123",
			Name:  "Jon Snow",
			Email: "jon.snow@got.com",
		}, nil
	}

	if provider == "google" && code == "456" {
		return &UserProfile{
			ID:    "GO456",
			Name:  "Bob",
			Email: "",
		}, nil
	}

	return nil, nil
}
