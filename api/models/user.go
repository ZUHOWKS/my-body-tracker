package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Age       int     `json:"age"`
	Weight    float64 `json:"weight"`
	Height    int     `json:"height"` // in cm
	Goal      string  `json:"goal"`
	Sex       int     `json:"sex"` // 0 for female, 1 for male
	Targets   Target  `json:"targets" gorm:"foreignKey:UserID"`
}

type Target struct {
	gorm.Model
	UserID   uint    `json:"userId" gorm:"uniqueIndex"`
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"` // in grams
	Carbs    float64 `json:"carbs"`   // in grams
	Fat      float64 `json:"fat"`     // in grams
	Fiber    float64 `json:"fiber"`   // in grams
}

type WeightRecord struct {
	gorm.Model
	UserID uint      `json:"userId" gorm:"index"`
	Weight float64   `json:"weight"`
	Date   time.Time `json:"date"`
	Note   string    `json:"note,omitempty"`
}
