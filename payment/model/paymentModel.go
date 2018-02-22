package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Payment struct {
	gorm.Model
	Uid           string     `gorm:"unique;not null"`
	AccountOrigin string     `gorm:"not null"`
	AccountTarget string     `gorm:"not null"`
	Amount        float64    `gorm:"not null"`
	Date          time.Time  `gorm:"not null"`
	Processed     bool       `gorm:"not null"`
	ProcessedDate *time.Time `gorm:"null"`
}

func SetUp(db *gorm.DB) *gorm.DB {

	db.SingularTable(true)
	db.AutoMigrate(&Payment{})

	return db
}
