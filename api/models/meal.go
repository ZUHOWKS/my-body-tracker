package models

import (
	"time"

	"gorm.io/gorm"
)

type MealType string

const (
	Breakfast MealType = "breakfast"
	Lunch     MealType = "lunch"
	Break     MealType = "break"
	Dinner    MealType = "dinner"
)

type Meal struct {
	gorm.Model
	Type   MealType  `json:"type" gorm:"column:meal_type;type:varchar(20)"`
	Date   time.Time `json:"date" gorm:"index"`
	UserID uint      `json:"userId" gorm:"column:user_id;index"`
	User   User      `json:"-" gorm:"foreignKey:UserID;references:ID"`
	Foods  []Food    `json:"foods" gorm:"many2many:meal_foods;"`
}
