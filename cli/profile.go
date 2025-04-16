package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func handleProfileCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: profile create | list | select <id> | view <id>")
		return
	}

	switch args[0] {
	case "create":
		user := promptUserProfile()
		createProfile(user)
	case "list":
		listProfiles()
	case "select":
		if len(args) != 2 {
			fmt.Println("Usage: profile select <id>")
			return
		}
		selectProfile(args[1])
	case "view":
		if len(args) != 2 {
			fmt.Println("Usage: profile view <id>")
			return
		}
		viewProfile(args[1])
	default:
		fmt.Println("Unknown profile command. Available: create, list, select, view")
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
	payload, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error creating profile:", err)
		return
	}

	resp, err := http.Post(apiURL+"/users", "application/json", bytes.NewBuffer(payload))
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

func listProfiles() {
	resp, err := http.Get(fmt.Sprintf("%s/users", apiURL))
	if err != nil {
		fmt.Println("Error listing users:", err)
		return
	}
	defer resp.Body.Close()

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	fmt.Println("Available users:")
	for _, user := range users {
		fmt.Printf("ID: %d - %s %s\n", user.ID, user.FirstName, user.LastName)
	}
}

func selectProfile(id string) {
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		fmt.Println("Invalid user ID")
		return
	}

	// Verify that the user exists
	resp, err := http.Get(fmt.Sprintf("%s/users/%s", apiURL, id))
	if err != nil {
		fmt.Println("Error selecting user:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("User not found")
		return
	}

	if err := setCurrentUserID(uint(userID)); err != nil {
		fmt.Println("Error saving session:", err)
		return
	}

	fmt.Printf("Selected user with ID %s\n", id)
}

func viewProfile(id string) {
	resp, err := http.Get(fmt.Sprintf("%s/users/%s", apiURL, id))
	if err != nil {
		fmt.Println("Error getting user profile:", err)
		return
	}
	defer resp.Body.Close()

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	fmt.Printf("Name: %s %s\nAge: %d\nHeight: %d cm\nWeight: %.2f kg\nGender: %s\n",
		user.FirstName, user.LastName, user.Age, user.Height, user.Weight, user.Gender)
}
