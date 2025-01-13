package database

import (
	"fmt"
	"log"

	"github.com/To-ge/gr_backend_go/config"
	"github.com/To-ge/gr_backend_go/infrastructure/database/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConnector struct {
	Conn *gorm.DB
}

func NewDBConnector() (*DBConnector, error) {
	conf := config.LoadConfig()
	dsn := combineDBInfo(*conf.DBInfo)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	models := []interface{}{model.User{}, model.Location{}, model.TelemetryLog{}}
	if err := db.AutoMigrate(models...); err != nil {
		log.Printf("Failed to migrate models: %s", err.Error())
		return nil, err
	}
	log.Println("completed auto migration.")
	return &DBConnector{Conn: db}, err
}

func combineDBInfo(info config.DBInfo) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Tokyo", info.Address, info.User, info.Password, info.DBName, info.DBPort)
}
