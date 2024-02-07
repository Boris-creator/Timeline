package oauth

import (
	"errors"
	"fiber-server/utils"
	"fmt"
	"net/url"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type GithubUserData struct {
	Login     string `json:"login"`
	GithubId  int    `json:"id"`
	AvatarUrl string `json:"avatar_url"`
}
type accessTokenResponse struct {
	AccessToken      string `json:"access_token"`
	Scope            string `json:"scope"`
	TokenType        string `json:"token_type"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type OAuthGH struct{}

func (OAuthGH) GetLoginUrl() url.URL {
	githubClientID := os.Getenv("GITHUB_CLIENT_ID")
	githubRedirectUri := "http://127.0.0.1:3000/github/callback"
	state, _ := bcrypt.GenerateFromPassword(GetGithubSecretForState(), 8)

	queryParams := url.Values{
		"client_id":    {githubClientID},
		"redirect_uri": {githubRedirectUri},
		"state":        {string(state)},
	}
	return url.URL{
		Scheme:   "https",
		Host:     "github.com",
		Path:     "/login/oauth/authorize",
		RawQuery: queryParams.Encode(),
	}
}

func (OAuthGH) GetAccessToken(code string) (string, error) {
	githubClientID := os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")

	requestBody := map[string]any{
		"client_id":     githubClientID,
		"client_secret": githubClientSecret,
		"code":          code,
	}

	responseData, err := utils.Fetch[accessTokenResponse](
		"https://github.com/login/oauth/access_token",
		"POST",
		&requestBody,
		utils.RequestOptions{
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
			},
		},
	)
	if err != nil {
		return "", err
	}
	if responseData.Error != "" {
		return "", errors.New(responseData.ErrorDescription)
	}

	return responseData.AccessToken, nil
}

func GetGithubData(accessToken string) (GithubUserData, error) {
	userData, err := utils.Fetch[GithubUserData](
		"https://api.github.com/user",
		"GET",
		nil,
		utils.RequestOptions{
			Headers: map[string]string{
				"Authorization": fmt.Sprintf("token %s", accessToken),
			},
		},
	)
	if err != nil {
		return *userData, err
	}

	return *userData, nil
}

func GetGithubSecretForState() []byte {
	return []byte(os.Getenv("GITHUB_STATE_SECRET"))
}
