package db

import (
	"fiber-server/models/event"
	"fiber-server/models/user"
	"fmt"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func ConnectDB() error {
	var err error
	env, err := godotenv.Read(".env")
	if err != nil {
		panic(err)
	}
	var DATABASE_URI string = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		env["DB_HOST"], env["DB_USER"], env["DB_PASSWORD"], env["DB_NAME"], env["DB_PORT"],
	)
	Database, err = gorm.Open(postgres.Open(DATABASE_URI), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		panic(err)
	}

	Database.AutoMigrate(
		&user.User{},
		&event.Event{},
		&event.EventDescription{},
		&event.EventRole{},
		&event.EventParticipant{},
	)

	return nil
}
