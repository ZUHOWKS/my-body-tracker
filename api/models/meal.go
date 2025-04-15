package models

import (
	"time"

	"gorm.io/gorm"
)

type MealPlan struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      uint   `json:"userId"`
	Meals       []Meal `json:"meals"`
}

type Meal struct {
	gorm.Model
	Name       string `json:"name"`
	MealPlanID uint   `json:"mealPlanId"`
	Foods      []Food `json:"foods" gorm:"many2many:meal_foods;"`
}

// ???
type DailyIntake struct {
	gorm.Model
	Date   time.Time `json:"date"`
	UserID uint      `json:"userId"`
	Meals  []Meal    `json:"meals" gorm:"many2many:daily_intake_meals;"`
}
