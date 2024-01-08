package event

import (
	"fiber-server/models/user"

	"gorm.io/gorm"
)

type EventParticipant struct {
	gorm.Model
	DateFrom      *string    `gorm:"type:time"`
	DateTo        *string    `gorm:"type:time"`
	PrecisionFrom *Precision `gorm:"type:varchar(50)"`
	PrecisionTo   *Precision `gorm:"type:varchar(50)"`
	EventId       uint       `gorm:"not null"`
	UserId        uint       `gorm:"not null"`
	RoleId        uint
	User          user.User
}

func (EventParticipant) TableName() string {
	return "participants"
}
