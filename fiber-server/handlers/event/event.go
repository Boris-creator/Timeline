package event

import (
	"fiber-server/auth"
	"fiber-server/db"
	"fiber-server/models/event"
	"sort"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type eventRole struct {
	Id                int                      `json:"id" validate:"omitempty,exists=roles"`
	Name              string                   `json:"name" validate:"required,max=255"`
	Description       string                   `json:"description" validate:"required"`
	EventParticipants []event.EventParticipant `json:"eventParticipants"`
}
type SaveEvent struct {
	Name          string      `json:"name" validate:"required,max=255"`
	Description   string      `json:"description" validate:"required"`
	DateFrom      string      `json:"dateFrom"`
	DateTo        string      `json:"dateTo"`
	PrecisionFrom string      `json:"precisionFrom" validate:"oneof=year month day hour"`
	PrecisionTo   string      `json:"precisionTo" validate:"oneof=year month day hour"`
	EventRoles    []eventRole `json:"eventRoles" validate:"required,dive"`
}

func StoreEvent(c *fiber.Ctx) error {
	var req = event.Event{}
	c.BodyParser(&req)

	eventParticipants := map[string][]event.EventParticipant{}
	for roleIdx, role := range req.EventRoles {
		eventParticipants[role.Name] = role.EventParticipants
		req.EventRoles[roleIdx].EventParticipants = []event.EventParticipant{}
	}

	userId := auth.GetCtxUserData(c).Id
	req.UserId = userId

	db.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&req).Error; err != nil {
			return err
		}

		for roleName, items := range eventParticipants {
			for _, participant := range items {
				roleIdx := sort.Search(len(req.EventRoles), func(i int) bool {
					return req.EventRoles[i].Name == roleName
				})
				participant.EventId = req.ID
				participant.UserId = userId
				participant.RoleId = req.EventRoles[roleIdx].ID
				if err := tx.Create(&participant).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})

	return c.JSON(req)
}

func UpdateEvent(c *fiber.Ctx) error {
	var req = event.Event{}
	c.BodyParser(&req)

	req.UserId = auth.GetCtxUserData(c).Id
	db.Database.Create(&req)

	return c.JSON(req)
}

func SearchEvents(c *fiber.Ctx) error {
	var events []event.Event
	db.Database.Preload("EventRoles").Find(&events)
	return c.JSON(events)
}
