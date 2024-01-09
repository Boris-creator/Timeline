package auth

import (
	"errors"
	"fiber-server/auth/oauth"
	"fiber-server/db"
	"fiber-server/models/user"
	"fiber-server/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RequestUser struct {
	Id    uint
	Login string
}

func Register(credentials Credentials) (user.User, error) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(credentials.Password), 8)
	user := user.User{
		Login:    credentials.Login,
		Password: string(hashed),
	}
	result := db.Database.Create(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func RegisterByGithub(account oauth.GithubUserData) (user.User, user.GithubUser) {
	githubUser := user.GithubUser{
		Login:     account.Login,
		GithubId:  account.GithubId,
		AvatarUrl: account.AvatarUrl,
	}

	var existingGithubUser user.GithubUser
	isFound := db.Database.Model(&githubUser).
		Preload("User").
		Where("github_id = ?", githubUser.GithubId).First(&existingGithubUser)
	if isFound.RowsAffected == 0 {
		user := user.User{
			Login:    githubUser.Login,
			Accounts: []user.GithubUser{githubUser},
		}
		db.Database.Create(&user)
		return user, githubUser
	}
	return existingGithubUser.User, existingGithubUser
}

func FindUserByCredentials(credentials Credentials) (user.User, error) {
	var existingUser user.User
	result := db.Database.Where(map[string]string{
		"login": credentials.Login,
	}).First(&existingUser)
	if result.RowsAffected == 0 {
		return existingUser, errors.New("")
	}
	err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(credentials.Password))
	if err != nil {
		return user.User{}, errors.New("")
	}
	return existingUser, nil
}

func GenerateToken(user user.User) string {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"login": user.Login,
		"exp":   time.Now().Add(time.Hour * 48).Unix(),
	}

	return utils.GenerateTokenString(claims)
}
func GenerateTokenWithGH(user user.User, account user.GithubUser) string {
	claims := jwt.MapClaims{
		"id":           user.ID,
		"login":        user.Login,
		"githubId":     account.GithubId,
		"githubAvatar": account.AvatarUrl,
		"exp":          time.Now().Add(time.Hour * 48).Unix(),
	}

	return utils.GenerateTokenString(claims)
}

func GetCtxUserData(c *fiber.Ctx) RequestUser {
	requestUser := c.Locals("user")
	currentUser, _ := requestUser.(jwt.MapClaims)

	return RequestUser{
		Id:    uint(currentUser["id"].(float64)),
		Login: currentUser["login"].(string),
	}
}
