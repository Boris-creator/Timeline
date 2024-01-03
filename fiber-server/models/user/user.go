package user

import (
	"fiber-server/models/event"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Password          string                   `gorm:"not null" json:"-"`
	Login             string                   `gorm:"not null;unique;"`
	EventDescriptions []event.EventDescription `gorm:"foreignKey:AuthorId"`
	Events            []event.Event            `gorm:"foreignKey:UserId"`
	EventParticipants []event.EventParticipant `gorm:"foreignKey:UserId;"`
}
