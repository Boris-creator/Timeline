package event

import (
	"fiber-server/auth"
	"fiber-server/db"
	"fiber-server/models/event"
	"sort"

	"github.com/gofiber/fiber/v2"
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
	PrecisionFrom string      `json:"precisionFrom"`
	PrecisionTo   string      `json:"precisionTo"`
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
	db.Database.Create(&req)

	for roleName, items := range eventParticipants {
		for _, participant := range items {
			roleIdx := sort.Search(len(req.EventRoles), func(i int) bool {
				return req.EventRoles[i].Name == roleName
			})
			participant.EventId = req.ID
			participant.UserId = userId
			participant.RoleId = req.EventRoles[roleIdx].ID
			db.Database.Create(&participant)
		}
	}

	return c.JSON(req)
}

func UpdateEvent(c *fiber.Ctx) error {
	var req = event.Event{}
	c.BodyParser(&req)

	req.UserId = auth.GetCtxUserData(c).Id
	db.Database.Create(&req)

	return c.JSON(req)
}
