package oauth

import (
	"fiber-server/auth"
	"fiber-server/auth/oauth"
	"net/url"
	"os"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func GithubLogin(c *fiber.Ctx) error {
	githubClientID := os.Getenv("GITHUB_CLIENT_ID")
	githubRedirectUri := "http://127.0.0.1:3000/github/callback"
	state, _ := bcrypt.GenerateFromPassword(getGithubSecretForState(), 8)

	queryParams := url.Values{
		"client_id":    {githubClientID},
		"redirect_uri": {githubRedirectUri},
		"state":        {string(state)},
	}
	redirectURL := url.URL{
		Scheme:   "https",
		Host:     "github.com",
		Path:     "/login/oauth/authorize",
		RawQuery: queryParams.Encode(),
	}

	return c.Status(fiber.StatusMovedPermanently).Redirect(redirectURL.String())
}

func GithubCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")
	stateError := bcrypt.CompareHashAndPassword([]byte(state), getGithubSecretForState())
	if stateError != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	githubAccessToken, err := oauth.GetGithubAccessToken(code)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).Send([]byte(err.Error()))
	}
	userData, _ := oauth.GetGithubData(githubAccessToken)
	user, account := auth.RegisterByGithub(userData)
	token := auth.GenerateTokenWithGH(user, account)
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		HTTPOnly: true,
	})
	return c.Redirect("/")
}

func getGithubSecretForState() []byte {
	return []byte(os.Getenv("GITHUB_STATE_SECRET"))
}
