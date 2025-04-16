package cli

type User struct {
	ID            uint    `json:"id"`
	FirstName     string  `json:"firstName"`
	LastName      string  `json:"lastName"`
	Name          string  `json:"name"`
	Age           int     `json:"age"`
	Weight        float64 `json:"weight"`
	Height        int     `json:"height"` // in cm
	Gender        string  `json:"gender"`
	Goal          string  `json:"goal"`
	Sex           int     `json:"sex"`           // 0 for female, 1 for male
	ActivityLevel int     `json:"activityLevel"` // 0-7 days of activity per week
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
	Protein  float64 `json:"protein"` // in grams
	Carbs    float64 `json:"carbs"`   // in grams
	Fat      float64 `json:"fat"`     // in grams
	Fiber    float64 `json:"fiber"`   // in grams
}
