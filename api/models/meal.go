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

type MealPlan struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      uint   `json:"userId" gorm:"column:user_id;index"`
	User        User   `json:"-" gorm:"foreignKey:UserID;references:ID"`
	Meals       []Meal `json:"meals" gorm:"foreignKey:MealPlanID;references:ID"`
}

type Meal struct {
	gorm.Model
	Type   MealType  `json:"type" gorm:"column:meal_type;type:varchar(20)"`
	Date   time.Time `json:"date" gorm:"index"`
	UserID uint      `json:"userId" gorm:"column:user_id;index"`
	User   User      `json:"-" gorm:"foreignKey:UserID;references:ID"`
	Foods  []Food    `json:"foods" gorm:"many2many:meal_foods;"`
}

// ???
type DailyIntake struct {
	gorm.Model
	Date   time.Time `json:"date" gorm:"index"`
	UserID uint      `json:"userId" gorm:"column:user_id;index"`
	User   User      `json:"-" gorm:"foreignKey:UserID;references:ID"`
	Meals  []Meal    `json:"meals" gorm:"many2many:daily_intake_meals;"`
}
