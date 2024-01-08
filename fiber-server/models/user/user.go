package user

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Password string       `json:"-"`
	Login    string       `gorm:"not null;unique;"`
	Accounts []GithubUser `gorm:"foreignKey:UserId;"`
}
