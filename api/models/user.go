package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName    string        `json:"firstName"`
	LastName     string        `json:"lastName"`
	Age          int          `json:"age"`
	Weight       float64      `json:"weight"`
	Height       int          `json:"height"` // in cm
	Goal         string       `json:"goal"`
	Sex          int          `json:"sex"` // 0 for female, 1 for male
	MealPlans    []MealPlan    `json:"mealPlans" gorm:"foreignKey:UserID"`
	DailyIntakes []DailyIntake `json:"dailyIntakes" gorm:"foreignKey:UserID"`
}
