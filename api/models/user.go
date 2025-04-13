package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName    string  `json:"firstName"`
	LastName     string  `json:"lastName"`
	Age          int     `json:"age"`
	Weight       float64 `json:"weight"`
	Height       int     `json:"height"` // in cm
	Goal         string  `json:"goal"`
	Sex          int     `json:"sex"` // 0 for female, 1 for male
	MealPlans    []MealPlan
	DailyIntakes []DailyIntake
}

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

type Food struct {
	gorm.Model
	FdcID       string  `json:"fdcId"`
	Name        string  `json:"name"`
	Protein     float64 `json:"protein"`
	Carbs       float64 `json:"carbs"`
	Fat         float64 `json:"fat"`
	Calories    float64 `json:"calories"`
	Fiber       float64 `json:"fiber"`
	ServingSize float64 `json:"servingSize"`
}

type DailyIntake struct {
	gorm.Model
	Date   time.Time `json:"date"`
	UserID uint      `json:"userId"`
	Meals  []Meal    `json:"meals" gorm:"many2many:daily_intake_meals;"`
}
