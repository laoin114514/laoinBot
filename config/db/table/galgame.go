package table

import (
	"gorm.io/gorm"
)

type Galgame struct {
	gorm.Model
	ID           int64  `gorm:"primaryKey;autoIncrement"`
	Title        string `gorm:"type:varchar(255);not null"`
	Url          string `gorm:"type:varchar(255);not null"`
	Introduction string `gorm:"type:text;not null"`
}

func (g *Galgame) TableName() string {
	return "galgame"
}
