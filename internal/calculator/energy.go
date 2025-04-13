package calculator

// CalculateBasalMetabolism calculates the Basal Metabolic Rate (BMR)
// weight in kg, height in meters, age in years, sex: 1 for male, 0 for female
func CalculateBasalMetabolism(weight float64, height int, age int, sex int) float64 {
	if weight <= 0 || height <= 0 || age <= 0 {
		return 0
	}

	coefficient := 0.963 // female coefficient
	if sex == 1 {
		coefficient = 1.083 // male coefficient
	}

	return (coefficient*weight*0.48*float64(height/100)*0.50*float64(age) - 0.13) * (1000 / 4.1855)
}

// CalculateActivityLevel calculates the Physical Activity Level (PAL)
// daysPerWeek: number of days of physical activity per week
func CalculateActivityLevel(daysPerWeek int) float64 {
	switch {
	case daysPerWeek <= 0:
		return 1.2 // Sedentary
	case daysPerWeek <= 3:
		return 1.375 // Light activity
	case daysPerWeek <= 5:
		return 1.55 // Moderate activity
	case daysPerWeek <= 6:
		return 1.725 // High activity
	default:
		return 1.9 // Very high activity
	}
}
