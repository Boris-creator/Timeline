package user

import (
	"gorm.io/gorm"
)

type GithubUser struct {
	gorm.Model
	Login     string `gorm:"not null;"`
	GithubId  int    `gorm:"not null;unique;"`
	AvatarUrl string `gorm:"text"`
	UserId    uint   `gorm:"not null"`
}
