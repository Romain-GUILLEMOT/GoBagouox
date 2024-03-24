package database

import (
	"GoBagouox/utils"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"os/signal"
	"syscall"
)

func StartConnexion() {
	username := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	dbname := os.Getenv("DATABASE_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		utils.Fatal("Can't make a connexion to the database!", err, 2)
	}

	err = db.Ping()
	if err != nil {
		utils.Fatal("Can't ping the database!", err, 2)
	}

	utils.Info("Database is now connected.", 2)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	utils.Info("Database has now disconnected.", 1)
	err = db.Close()
	if err != nil {
		utils.Fatal("Can't close the connexion to the database!", err, 2)
	}
}
