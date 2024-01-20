package oauth

import (
	"fiber-server/auth"
	"fiber-server/auth/oauth"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func GithubLogin(c *fiber.Ctx) error {
	oauthGH := oauth.OAuthGH{}
	redirectURL := oauthGH.GetLoginUrl()
	return c.Status(fiber.StatusMovedPermanently).Redirect(redirectURL.String())
}

func GithubCallback(c *fiber.Ctx) error {
	oauthGH := oauth.OAuthGH{}
	code := c.Query("code")
	state := c.Query("state")
	stateError := bcrypt.CompareHashAndPassword([]byte(state), oauth.GetGithubSecretForState())
	if stateError != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	githubAccessToken, err := oauthGH.GetAccessToken(code)
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
