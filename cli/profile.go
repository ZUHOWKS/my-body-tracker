package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

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

	fmt.Print("Height (cm): ")
	fmt.Scanf("%d\n", &user.Height)

	fmt.Print("Sex (0 for female, 1 for male): ")
	fmt.Scanf("%d\n", &user.Sex)

	fmt.Print("Activity Level (0-7 days per week): ")
	for {
		fmt.Scanf("%d\n", &user.ActivityLevel)
		if user.ActivityLevel >= 0 && user.ActivityLevel <= 7 {
			break
		}
		fmt.Print("Please enter a number between 0 and 7: ")
	}

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

	fmt.Printf("Height: %d cm\n", int(data["height"].(float64)))
	fmt.Printf("Weight: %.2f kg\n", data["weight"])
	fmt.Printf("BMI: %.2f\n", data["bmi"])
	fmt.Printf("Body Fat Percentage: %.2f%%\n", data["bfp"])
	fmt.Printf("Basal Metabolic Rate: %.2f kcal/day\n", data["bmr"])
}
