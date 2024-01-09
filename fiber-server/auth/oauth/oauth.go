package oauth

import (
	"errors"
	"fiber-server/utils"
	"fmt"
	"os"
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

func GetGithubAccessToken(code string) (string, error) {
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
