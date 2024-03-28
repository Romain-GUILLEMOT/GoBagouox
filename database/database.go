package database

import (
	"GoBagouox/database/models"
	"GoBagouox/utils"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
)

var db *gorm.DB

type ModelWithName struct {
	Model interface{}
	Name  string
}

func getModals() []interface{} {
	return []interface{}{
		ModelWithName{Model: &models.User{}, Name: "Users"},
		ModelWithName{Model: &models.Ticket{}, Name: "Tickets"},
		ModelWithName{Model: &models.TicketMessage{}, Name: "Ticket_messages"},
		ModelWithName{Model: &models.TicketAttachments{}, Name: "Ticket_attachments"},
		ModelWithName{Model: &models.SshAccess{}, Name: "Sshaccess"},
	}
}
func Migrate(modelWithNames ...interface{}) {
	for _, modelWithName := range modelWithNames {
		mwn := modelWithName.(ModelWithName) // This will panic if the assertion is not ok
		err := db.AutoMigrate(mwn.Model)
		if err != nil {
			utils.Fatal("Can't run all models.", err, 2)
		}
		utils.Debug("Model "+utils.Blue(mwn.Name)+" was migrated!", 2)
	}
}

func GetDB() *gorm.DB {
	return db
}

func StartConnexion() {
	username := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	dbname := os.Getenv("DATABASE_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, dbname)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.Fatal("Can't make a connexion to the database!", err, 2)
	}
	utils.Info("Database is now connected.", 2)
	models := getModals()
	Migrate(models...)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	utils.Info("Database has now disconnected.", 1)
}
