package oauth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

	requestBody := map[string]string{
		"client_id":     githubClientID,
		"client_secret": githubClientSecret,
		"code":          code,
	}
	requestJSON, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	responseBody, _ := io.ReadAll(response.Body)
	var responseData accessTokenResponse
	json.Unmarshal(responseBody, &responseData)
	if responseData.Error != "" {
		return "", errors.New(responseData.ErrorDescription)
	}

	return responseData.AccessToken, nil
}

func GetGithubData(accessToken string) (GithubUserData, error) {
	var userData GithubUserData

	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		return userData, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return userData, err
	}
	responseBody, _ := io.ReadAll(response.Body)
	json.Unmarshal(responseBody, &userData)
	return userData, nil
}
