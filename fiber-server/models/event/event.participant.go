package event

import "gorm.io/gorm"

type EventParticipant struct {
	gorm.Model
	DateFrom      *string `gorm:"type:time"`
	DateTo        *string `gorm:"type:time"`
	PrecisionFrom *string `gorm:"type:varchar(50)"`
	PrecisionTo   *string `gorm:"type:varchar(50)"`
	EventId       uint    `gorm:"not null"`
	UserId        uint    `gorm:"not null"`
	RoleId        uint
}

func (EventParticipant) TableName() string {
	return "participants"
}
