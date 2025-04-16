package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func handleMealCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: meal add <type> <date> <food_name> | view <date> <type> | list [date] [type]")
		return
	}

	switch args[0] {
	case "add":
		if len(args) != 4 {
			fmt.Println("Usage: meal add <type> <date> <food_name>")
			fmt.Println("  type: breakfast, lunch, break, dinner")
			fmt.Println("  date: YYYY-MM-DD or 'today'")
			fmt.Println("  food_name: name of the food to search for")
			return
		}
		addFoodToMealType(args[1], args[2], args[3])
	case "view":
		if len(args) != 3 {
			fmt.Println("Usage: meal view <type> <date>")
			fmt.Println("  date: YYYY-MM-DD or 'today'")
			fmt.Println("  type: breakfast, lunch, break, dinner")
			return
		}
		viewMeal(args[1], args[2])
	case "list":
		userID := getCurrentUserID()
		if userID == 0 {
			fmt.Println("No user selected. Use 'profile select <id>' to select one.")
			return
		}
		var date, mealType string
		if len(args) > 1 {
			mealType = args[1]
		}
		if len(args) > 2 {
			date = args[2]
		}
		listMeals(fmt.Sprint(userID), date, mealType)
	default:
		fmt.Println("Unknown meal command")
	}
}

func viewMeal(mealType, dateStr string) {
	// Check if a user is selected
	userID := getCurrentUserID()
	if userID == 0 {
		fmt.Println("No user selected. Use 'profile list' to see available users and 'profile select <id>' to select one.")
		return
	}

	// Parse date
	var date time.Time
	if dateStr == "today" {
		date = time.Now()
	} else {
		var err error
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			fmt.Println("Invalid date format. Use YYYY-MM-DD or 'today'")
			return
		}
	}

	// Validate meal type
	mealType = strings.ToLower(mealType)
	validTypes := map[string]bool{"breakfast": true, "lunch": true, "break": true, "dinner": true}
	if !validTypes[mealType] {
		fmt.Println("Invalid meal type. Must be one of: breakfast, lunch, break, dinner")
		return
	}

	// Get meal details
	resp, err := http.Get(fmt.Sprintf("%s/meals/user/%d?date=%s&type=%s", apiURL, userID, date.Format("2006-01-02"), mealType))
	if err != nil {
		fmt.Println("Error getting meal:", err)
		return
	}
	defer resp.Body.Close()

	var meals []struct {
		ID    uint      `json:"id"`
		Name  string    `json:"name"`
		Type  string    `json:"type"`
		Date  time.Time `json:"date"`
		Foods []struct {
			ID          string  `json:"id"`
			Name        string  `json:"name"`
			Calories    float64 `json:"calories"`
			Protein     float64 `json:"protein"`
			Carbs       float64 `json:"carbs"`
			Fat         float64 `json:"fat"`
			Fiber       float64 `json:"fiber"`
			Sugar       float64 `json:"sugar"`
			ServingSize float64 `json:"servingSize"`
		} `json:"foods"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&meals); err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	if len(meals) == 0 {
		fmt.Printf("No %s found for %s\n", mealType, date.Format("2006-01-02"))
		return
	}

	meal := meals[0]
	fmt.Printf("\n%s - %s\n", meal.Type, meal.Date.Format("2006-01-02"))
	fmt.Println("Foods:")

	// Track total nutrients
	var totalNutrients struct {
		Calories float64
		Protein  float64
		Carbs    float64
		Fat      float64
		Fiber    float64
		Sugar    float64
	}

	// List foods and accumulate nutrients
	for _, food := range meal.Foods {
		fmt.Printf("- %s (%.0f calories, %.1fg protein, %.1fg carbs, %.1fg fat, %.1fg fiber, %.1fg sugar)\n", food.Name, food.Calories, food.Protein, food.Carbs, food.Fat, food.Fiber, food.Sugar)
		totalNutrients.Calories += food.Calories
		totalNutrients.Protein += food.Protein
		totalNutrients.Carbs += food.Carbs
		totalNutrients.Fat += food.Fat
		totalNutrients.Fiber += food.Fiber
		totalNutrients.Sugar += food.Sugar
	}

	// Display total nutrients
	fmt.Println("\nTotal Nutrients:")
	fmt.Printf("Calories: %.0f\n", totalNutrients.Calories)
	fmt.Printf("Protein: %.1fg\n", totalNutrients.Protein)
	fmt.Printf("Carbs: %.1fg\n", totalNutrients.Carbs)
	fmt.Printf("Fat: %.1fg\n", totalNutrients.Fat)
	fmt.Printf("Fiber: %.1fg\n", totalNutrients.Fiber)
	fmt.Printf("Sugar: %.1fg\n", totalNutrients.Sugar)
}

func listMeals(userID, date, mealType string) {
	url := fmt.Sprintf("%s/meals/user/%s", apiURL, userID)
	if date != "" {
		url += "?date=" + date
	}
	if mealType != "" {
		if strings.Contains(url, "?") {
			url += "&"
		} else {
			url += "?"
		}
		url += "type=" + mealType
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error getting meals:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var message struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&message); err != nil {
			fmt.Println("Error parsing error response:", err)
		} else {
			fmt.Printf("Error: Server returned %s : %s\n", resp.Status, message.Error)
		}
		return
	}

	var meals []struct {
		ID    uint      `json:"id"`
		Name  string    `json:"name"`
		Type  string    `json:"type"`
		Date  time.Time `json:"date"`
		Foods []struct {
			Name     string  `json:"name"`
			Calories float64 `json:"calories"`
		} `json:"foods"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&meals); err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	for _, meal := range meals {
		fmt.Printf("\nMeal ID: %d\nName: %s\nType: %s\nDate: %s\n",
			meal.ID, meal.Name, meal.Type, meal.Date.Format("2006-01-02"))
		if len(meal.Foods) > 0 {
			fmt.Println("Foods:")
			for _, food := range meal.Foods {
				fmt.Printf("  - %s (%.0f calories)\n", food.Name, food.Calories)
			}
		}
	}
}

func addFoodToMealType(mealType, dateStr, foodQuery string) {
	// Check if a user is selected
	userID := getCurrentUserID()
	if userID == 0 {
		fmt.Println("No user selected. Use 'profile list' to see available users and 'profile select <id>' to select one.")
		return
	}

	// Validate meal type
	mealType = strings.ToLower(mealType)
	validTypes := map[string]bool{"breakfast": true, "lunch": true, "break": true, "dinner": true}
	if !validTypes[mealType] {
		fmt.Println("Invalid meal type. Must be one of: breakfast, lunch, break, dinner")
		return
	}

	// Parse date
	var date time.Time
	if dateStr == "today" {
		date = time.Now()
	} else {
		var err error
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			fmt.Println("Invalid date format. Use YYYY-MM-DD or 'today'")
			return
		}
	}

	// Search for food
	resp, err := http.Get(fmt.Sprintf("%s/foods/search?q=%s", apiURL, url.QueryEscape(foodQuery)))
	if err != nil {
		fmt.Println("Error searching for food:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Define food item structure
	type FoodItem struct {
		FdcID    string  `json:"fdcId"`
		Name     string  `json:"name"`
		Calories float64 `json:"calories"`
	}

	// First, try to parse the raw JSON to see what we're dealing with
	var rawData map[string]interface{}
	if err := json.Unmarshal(respBody, &rawData); err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	// Extract foods from the response
	var foods []FoodItem

	// Check if the response has a 'foods' field
	if foodsData, ok := rawData["foods"]; ok {
		// Convert the foods data to JSON
		foodsJSON, err := json.Marshal(foodsData)
		if err != nil {
			fmt.Println("Error marshaling foods data:", err)
			return
		}

		// Unmarshal into our foods slice
		if err := json.Unmarshal(foodsJSON, &foods); err != nil {
			fmt.Println("Error parsing foods data:", err)
			return
		}
	} else {
		// Try parsing the whole response as an array of foods
		if err := json.Unmarshal(respBody, &foods); err != nil {
			// If that fails, try with a results field
			var resultsResponse struct {
				Results []FoodItem `json:"results"`
			}
			if err := json.Unmarshal(respBody, &resultsResponse); err != nil {
				fmt.Println("Error parsing food search results:", err)
				return
			}
			foods = resultsResponse.Results
		}
	}

	if len(foods) == 0 {
		fmt.Println("No foods found matching your query")
		return
	}

	// Display food options
	fmt.Println("\nFound foods:")
	for i, food := range foods {
		fmt.Printf("%d. %s (%.0f calories)\n", i+1, food.Name, food.Calories)
	}

	// Get user selection
	fmt.Print("\nSelect a food (enter number): ")
	var selection int
	fmt.Scanf("%d\n", &selection)
	selection-- // Convert to 0-based index

	if selection < 0 || selection >= len(foods) {
		fmt.Println("Invalid selection")
		return
	}

	// Create or get meal for the given type and date
	meal := struct {
		ID     uint      `json:"id"`
		Name   string    `json:"name"`
		Type   string    `json:"type"`
		Date   time.Time `json:"date"`
		UserID uint      `json:"userId"`
	}{
		Name:   fmt.Sprintf("%s on %s", mealType, date.Format("2006-01-02")),
		Type:   mealType,
		Date:   date,
		UserID: userID,
	}

	// Try to find existing meal
	mealURL := fmt.Sprintf("%s/meals/user/%d?date=%s&type=%s", apiURL, userID, date.Format("2006-01-02"), mealType)
	fmt.Println(mealURL)
	resp, err = http.Get(mealURL)
	if err != nil {
		fmt.Println("Error checking for existing meal:", err)
		return
	}
	defer resp.Body.Close()

	// Read the meal check response
	mealRespBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading meal check response:", err)
		return
	}

	// Initialize mealID
	var mealID uint

	// First check if we got an error response
	var errorResp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(mealRespBody, &errorResp); err == nil && errorResp.Error != "" {
		fmt.Println("Server returned an error:", errorResp.Error)
		// We'll create a new meal since we couldn't find an existing one
	} else {
		// Try to parse as an array of meals
		var existingMeals []struct {
			ID uint `json:"id"`
		}
		if err := json.Unmarshal(mealRespBody, &existingMeals); err == nil && len(existingMeals) > 0 {
			mealID = existingMeals[0].ID
		}
	}

	// If no existing meal was found, create a new one
	if mealID == 0 {
		// Create new meal
		payload, err := json.Marshal(meal)
		if err != nil {
			fmt.Println("Error preparing meal creation request:", err)
			return
		}

		resp, err = http.Post(apiURL+"/meals", "application/json", bytes.NewBuffer(payload))
		if err != nil {
			fmt.Println("Error creating meal:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			fmt.Println("Error: Server returned", resp.Status)
			return
		}

		var createdMeal struct {
			ID uint `json:"id"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&createdMeal); err != nil {
			fmt.Println("Error parsing meal creation response:", err)
			return
		}
		mealID = createdMeal.ID
	}

	// Add food to meal
	addFoodURL := fmt.Sprintf("%s/meals/%d/foods", apiURL, mealID)
	payload := map[string]string{"foodId": foods[selection].FdcID}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error preparing food addition request:", err)
		return
	}

	resp, err = http.Post(addFoodURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error adding food to meal:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var message struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&message); err != nil {
			fmt.Println("Error parsing error response:", err)
		} else {
			fmt.Printf("Error: Server returned %s : %s\n", resp.Status, message.Error)
		}
		return
	}

	fmt.Printf("Added %s to your %s for %s\n", foods[selection].Name, mealType, date.Format("2006-01-02"))
}
