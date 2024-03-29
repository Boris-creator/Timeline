package event

import (
	"fiber-server/models/user"

	"gorm.io/gorm"
)

type Precision string

const (
	Year  Precision = "year"
	Month Precision = "month"
	Day   Precision = "day"
	Hour  Precision = "hour"
)

type Event struct {
	gorm.Model
	ID                uint               `gorm:"primarykey" json:"id"`
	Name              string             `gorm:"type:varchar(255);not null;" json:"name"`
	Description       string             `gorm:"type:text;" json:"description"`
	DateFrom          string             `gorm:"type:time;" json:"dateFrom"`
	DateTo            string             `gorm:"type:time;" json:"dateTo"`
	PrecisionFrom     Precision          `gorm:"type:varchar(50);"`
	PrecisionTo       Precision          `gorm:"type:varchar(50);"`
	UserId            uint               `gorm:"not null"`
	Descriptions      []EventDescription `gorm:"foreignKey:EventId;"`
	EventRoles        []EventRole        `gorm:"many2many:event_roles;" json:"eventRoles"`
	EventParticipants []EventParticipant `gorm:"foreignKey:EventId;"`
	User              user.User
}
