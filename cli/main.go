package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const apiURL = "http://localhost:8080"

type User struct {
	ID        uint    `json:"id"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Age       int     `json:"age"`
	Weight    float64 `json:"weight"`
	Height    float64 `json:"height"`
	Goal      string  `json:"goal"`
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

func main() {
	fmt.Println("Welcome to My Body Tracker CLI!")
	fmt.Println("Available commands:")
	fmt.Println("  profile create - Create a new user profile")
	fmt.Println("  profile view <id> - View user profile and stats")
	fmt.Println("  food search <query> - Search for food in database")
	fmt.Println("  meal create - Create a new meal")
	fmt.Println("  meal add <meal_id> <food_id> - Add food to meal")
	fmt.Println("  meal view <id> - View meal details")
	fmt.Println("  exit")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if input == "exit" {
			break
		}

		handleCommand(input)
	}
}

func handleCommand(input string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	command := parts[0]
	args := parts[1:]

	switch command {
	case "profile":
		if len(args) == 0 {
			fmt.Println("Usage: profile <create|view> [id]")
			return
		}
		handleProfileCommand(args)

	case "food":
		if len(args) < 2 {
			fmt.Println("Usage: food search <query>")
			return
		}
		if args[0] == "search" {
			searchFood(strings.Join(args[1:], " "))
		}

	case "meal":
		if len(args) == 0 {
			fmt.Println("Usage: meal <create|add|view> [args...]")
			return
		}
		handleMealCommand(args)

	default:
		fmt.Println("Unknown command. Type 'help' for available commands")
	}
}

func handleProfileCommand(args []string) {
	switch args[0] {
	case "create":
		user := promptUserProfile()
		createProfile(user)
	case "view":
		if len(args) != 2 {
			fmt.Println("Usage: profile view <id>")
			return
		}
		viewProfile(args[1])
	default:
		fmt.Println("Unknown profile command. Available: create, view")
	}
}

func promptUserProfile() User {
	reader := bufio.NewReader(os.Stdin)
	var user User

	fmt.Print("First Name: ")
	user.FirstName, _ = reader.ReadString('\n')
	user.FirstName = strings.TrimSpace(user.FirstName)

	fmt.Print("Last Name: ")
	user.LastName, _ = reader.ReadString('\n')
	user.LastName = strings.TrimSpace(user.LastName)

	fmt.Print("Age: ")
	fmt.Scanf("%d\n", &user.Age)

	fmt.Print("Weight (kg): ")
	fmt.Scanf("%f\n", &user.Weight)

	fmt.Print("Height (m): ")
	fmt.Scanf("%f\n", &user.Height)

	fmt.Print("Goal (e.g., weight loss, muscle gain): ")
	user.Goal, _ = reader.ReadString('\n')
	user.Goal = strings.TrimSpace(user.Goal)

	return user
}

func createProfile(user User) {
	jsonData, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error preparing data:", err)
		return
	}

	resp, err := http.Post(apiURL+"/users", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating profile:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Println("Error: Server returned", resp.Status)
		return
	}

	fmt.Println("Profile created successfully!")
}

func viewProfile(id string) {
	resp, err := http.Get(apiURL + "/users/" + id)
	if err != nil {
		fmt.Println("Error fetching profile:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Server returned", resp.Status)
		return
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	fmt.Printf("Height: %.2f m\n", data["height"])
	fmt.Printf("Weight: %.2f kg\n", data["weight"])
	fmt.Printf("BMI: %.2f\n", data["bmi"])
	fmt.Printf("Body Fat Percentage: %.2f%%\n", data["bfp"])
	fmt.Printf("Basal Metabolic Rate: %.2f kcal/day\n", data["bmr"])
}

func searchFood(query string) {
	resp, err := http.Get(apiURL + "/foods/search?q=" + query)
	if err != nil {
		fmt.Println("Error searching food:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Server returned", resp.Status)
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	foods := result["foods"].([]interface{})
	for _, food := range foods {
		f := food.(map[string]interface{})
		fmt.Printf("ID: %s\nName: %s\n\n", f["fdcId"], f["description"])
	}
}

func handleMealCommand(args []string) {
	switch args[0] {
	case "create":
		// TODO: Implement meal creation
		fmt.Println("Meal creation not implemented yet")
	case "add":
		if len(args) != 3 {
			fmt.Println("Usage: meal add <meal_id> <food_id>")
			return
		}
		// TODO: Implement adding food to meal
		fmt.Println("Adding food to meal not implemented yet")
	case "view":
		if len(args) != 2 {
			fmt.Println("Usage: meal view <id>")
			return
		}
		// TODO: Implement meal viewing
		fmt.Println("Meal viewing not implemented yet")
	default:
		fmt.Println("Unknown meal command. Available: create, add, view")
	}
}
