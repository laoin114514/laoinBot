package table

import (
	"gorm.io/gorm"
)

type Galgame struct {
	gorm.Model
	ID           uint   `gorm:"primaryKey;type:int4;column:id"`
	Title        string `gorm:"type:varchar(255);not null;column:title"`
	Url          string `gorm:"type:varchar(255);not null;column:url"`
	Introduction string `gorm:"type:text;column:introduction"`
	Code         string `gorm:"type:varchar(50);column:code"`
	Type         string `gorm:"type:varchar(50);column:type"`
	Tag          string `gorm:"type:varchar(255);column:tag"`
	Name         string `gorm:"type:varchar(255);column:name"`
	Platform     string `gorm:"type:varchar(50);column:platform"`
	Size         string `gorm:"type:varchar(20);column:size"`
	Version      string `gorm:"type:varchar(100);column:version"`
}

func (g *Galgame) TableName() string {
	return "galgame"
}
