package cli

import (
	"encoding/json"
	"fmt"
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

	var result struct {
		Foods []Food `json:"foods"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	for _, food := range result.Foods {
		fmt.Printf("ID: %s\nName: %s\n\n", food.FdcID, food.Name)
	}
}
