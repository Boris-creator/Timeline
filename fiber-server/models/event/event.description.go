package event

import "gorm.io/gorm"

type EventDescription struct {
	gorm.Model
	Content  string `gorm:"type:text;not null;"`
	AuthorId uint   `gorm:"not null;"`
	EventId  uint   `gorm:"not null;"`
}

func (EventDescription) TableName() string {
	return "descriptions"
}
