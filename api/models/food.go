package models

import "gorm.io/gorm"

type Food struct {
	gorm.Model
	FdcID       string  `json:"fdcId" gorm:"uniqueIndex"`
	Name        string  `json:"name"`
	Protein     float64 `json:"protein"`
	Carbs       float64 `json:"carbs"`
	Fat         float64 `json:"fat"`
	Calories    float64 `json:"calories"`
	Fiber       float64 `json:"fiber"`
	ServingSize float64 `json:"servingSize"`
	Meals       []Meal  `json:"meals" gorm:"many2many:meal_foods;"`
}
