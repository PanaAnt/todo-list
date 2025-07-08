package database

import (
	"fmt"
	"log"
	"todoApp/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB() {
	// create connection str
	connStr := "user=postgres password=password dbname=todoApp port=5432 sslmode=disable"

	// open db connection
	var err error
	Db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	//migrate the schema
	err = Db.AutoMigrate(&models.User{}, &models.Todo{})
	if err != nil {
		log.Fatal("failed to migrate the schema", err)
	}

	fmt.Println("connected to db successfully!")
}
