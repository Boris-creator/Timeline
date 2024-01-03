package event

import "gorm.io/gorm"

type EventRole struct {
	gorm.Model
	Name              string             `gorm:"not null;"`
	Description       string             `gorm:"type:text;"`
	EventParticipants []EventParticipant `gorm:"foreignKey:RoleId;"`
}

func (EventRole) TableName() string {
	return "roles"
}
