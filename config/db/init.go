package db

import (
	"fmt"
	"laoinBot/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDb() error {
	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.BotConfig.Db.Host, config.BotConfig.Db.Port, config.BotConfig.Db.User, config.BotConfig.Db.Password, config.BotConfig.Db.Dbname)), &gorm.Config{})
	if err != nil {
		return err
	}
	Db = db
	return nil
}
