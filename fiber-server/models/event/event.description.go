package event

import (
	"fiber-server/models/user"

	"gorm.io/gorm"
)

type EventDescription struct {
	gorm.Model
	Content  string    `gorm:"type:text;not null;"`
	AuthorId uint      `gorm:"not null;"`
	EventId  uint      `gorm:"not null;"`
	User     user.User `gorm:"foreignKey:AuthorId;"`
}

func (EventDescription) TableName() string {
	return "descriptions"
}
