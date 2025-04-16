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
	Height    int     `json:"height"`
	Goal      string  `json:"goal"`
	Sex       int     `json:"sex"`
	Targets   Target  `json:"targets" gorm:"foreignKey:UserID"`
}

type Target struct {
	gorm.Model
	UserID   uint    `json:"userId" gorm:"uniqueIndex"`
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
	Fiber    float64 `json:"fiber"`
}

type WeightRecord struct {
	gorm.Model
	UserID uint      `json:"userId" gorm:"index"`
	Weight float64   `json:"weight"`
	Date   time.Time `json:"date"`
	Note   string    `json:"note,omitempty"`
}
