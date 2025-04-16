package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func handleFoodCommand(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: food search <query>")
		return
	}
	if args[0] == "search" {
		searchFood(strings.Join(args[1:], " "))
	}
}

func searchFood(query string) {
	fmt.Println(apiURL + "/foods/search?q=" + query)
	resp, err := http.Get(apiURL + "/foods/search?q=" + query)
	if err != nil {
		fmt.Println("Error searching food:", err)
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

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// First, try to parse the raw JSON to see what we're dealing with
	var rawData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &rawData); err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	// Debug: print the raw data structure
	fmt.Println("Response structure:", rawData)

	// Extract foods from the response
	var foods []Food
	
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
		if err := json.Unmarshal(bodyBytes, &foods); err != nil {
			fmt.Println("Could not parse response as foods array:", err)
			return
		}
	}

	if len(foods) == 0 {
		fmt.Println("No foods found matching your query")
		return
	}

	fmt.Printf("Found %d foods matching your query:\n\n", len(foods))
	for _, food := range foods {
		fmt.Printf("ID: %s\nName: %s\nCalories: %.0f\n\n", food.FdcID, food.Name, food.Calories)
	}
}
