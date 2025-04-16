package calculator

// CalculateBMI calculates the Body Mass Index (BMI)
// weight in kg, height in meters
func CalculateBMI(weight float64, height int) float64 {
	if height == 0 {
		return 0
	}
	return weight / ((float64(height) / 100) * (float64(height) / 100))
}

// CalculateBFP calculates the Body Fat Percentage (BFP)
// age in years, sex: 1 for male, 0 for female
func CalculateBFP(bmi float64, age int, sex int) float64 {
	if bmi <= 0 || age <= 0 || (sex != 0 && sex != 1) {
		return 0
	}
	return (1.20 * bmi) + (0.23 * float64(age)) - (10.8 * float64(sex)) - 5.4
}

// CalculateIMG calculates the Body Fat Mass Index (IMG)
// weight in kg, height in meters, sex: 1 for male, 0 for female
func CalculateIMG(weight float64, height int, sex int) float64 {
	if weight <= 0 || height <= 0 || (sex != 0 && sex != 1) {
		return 0
	}

	// Calculate body fat percentage using the formula
	// For men: (1.20 × BMI) + (0.23 × Age) - 16.2
	// For women: (1.20 × BMI) + (0.23 × Age) - 5.4
	// Since we don't have age in this function, we'll use a simplified version

	bmi := CalculateBMI(weight, height)

	// Simplified IMG calculation
	// IMG = (1.20 × BMI) - (10.8 × Sex) - 5.4
	return (1.20 * bmi) - (10.8 * float64(sex)) - 5.4
}
