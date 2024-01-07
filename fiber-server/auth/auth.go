package auth

import (
	"errors"
	"fiber-server/auth/oauth"
	"fiber-server/db"
	"fiber-server/models/user"
	"os"
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

func RegisterByGithub(account oauth.GithubUserData) {
	githubUser := user.GithubUser{
		Login:     account.Login,
		GithubId:  account.GithubId,
		AvatarUrl: account.AvatarUrl,
	}
	var count int64
	db.Database.Model(&githubUser).Where("github_id = ?", githubUser.GithubId).Count(&count)
	if count == 0 {
		user := user.User{
			Login:    githubUser.Login,
			Accounts: []user.GithubUser{githubUser},
		}
		db.Database.Create(&user)
	}
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return t
}

func GetCtxUserData(c *fiber.Ctx) RequestUser {
	requestUser := c.Locals("user")
	currentUser, _ := requestUser.(jwt.MapClaims)

	return RequestUser{
		Id:    uint(currentUser["id"].(float64)),
		Login: currentUser["login"].(string),
	}
}
