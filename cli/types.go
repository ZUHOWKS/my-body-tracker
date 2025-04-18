package cli

import (
	"time"
)

type User struct {
	ID            uint    `json:"id"`
	FirstName     string  `json:"firstName"`
	LastName      string  `json:"lastName"`
	Name          string  `json:"name"`
	Age           int     `json:"age"`
	Weight        float64 `json:"weight"`
	Height        int     `json:"height"`
	Gender        string  `json:"gender"`
	Goal          string  `json:"goal"`
	Sex           int     `json:"sex"`
	ActivityLevel int     `json:"activityLevel"`
}

type Food struct {
	FdcID       string  `json:"fdcId"`
	Name        string  `json:"name"`
	Protein     float64 `json:"protein"`
	Carbs       float64 `json:"carbs"`
	Fat         float64 `json:"fat"`
	Calories    float64 `json:"calories"`
	Fiber       float64 `json:"fiber"`
	ServingSize float64 `json:"servingSize"`
}

type Target struct {
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
	Fiber    float64 `json:"fiber"`
}

type WeightRecord struct {
	ID     uint      `json:"id"`
	Weight float64   `json:"weight"`
	Date   time.Time `json:"date"`
	Note   string    `json:"note,omitempty"`
}
