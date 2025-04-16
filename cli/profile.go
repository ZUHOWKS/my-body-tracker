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
	"time"
)

func handleProfileCommand(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: profile create | list | select <id> | view <id> | targets <id> | set-targets <id> | weight <id> | weight-history <id>")
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
	case "targets":
		if len(args) != 2 {
			fmt.Println("Usage: profile targets <id>")
			return
		}
		viewTargets(args[1])
	case "set-targets":
		if len(args) != 2 {
			fmt.Println("Usage: profile set-targets <id>")
			return
		}
		setTargets(args[1])
	case "weight":
		if len(args) != 2 {
			fmt.Println("Usage: profile weight <id>")
			return
		}
		recordWeight(args[1])
	case "weight-history":
		if len(args) != 2 {
			fmt.Println("Usage: profile weight-history <id>")
			return
		}
		viewWeightHistory(args[1])
	default:
		fmt.Println("Unknown profile command. Available: create, list, select, view, targets, set-targets, weight, weight-history")
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

	gender := "Female"
	if user.Sex == 1 {
		gender = "Male"
	}

	fmt.Printf("Name: %s %s\nAge: %d\nHeight: %d cm\nWeight: %.2f kg\nGender: %s\n",
		user.FirstName, user.LastName, user.Age, user.Height, user.Weight, gender)

	// Afficher les statistiques de l'utilisateur
	viewUserStats(id)
}

func viewUserStats(id string) {
	resp, err := http.Get(fmt.Sprintf("%s/users/%s/stats", apiURL, id))
	if err != nil {
		fmt.Println("Error getting user stats:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Server returned", resp.Status)
		return
	}

	var stats struct {
		Height int     `json:"height"`
		Weight float64 `json:"weight"`
		BMI    float64 `json:"bmi"`
		BFP    float64 `json:"bfp"`
		IMG    float64 `json:"img"`
		BMR    float64 `json:"bmr"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		fmt.Println("Error parsing stats response:", err)
		return
	}

	fmt.Println("\nHealth Statistics:")
	fmt.Printf("BMI: %.2f\n", stats.BMI)
	fmt.Printf("Body Fat Percentage: %.2f%%\n", stats.BFP)
	fmt.Printf("Body Fat Mass Index: %.2f\n", stats.IMG)
	fmt.Printf("Basal Metabolic Rate: %.2f kcal/day\n", stats.BMR)
}

func viewTargets(id string) {
	// Récupérer les objectifs
	resp, err := http.Get(fmt.Sprintf("%s/users/%s/targets", apiURL, id))
	if err != nil {
		fmt.Println("Error getting user targets:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		fmt.Println("No targets found for this user. Use 'profile set-targets <id>' to set targets.")
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Server returned", resp.Status)
		return
	}

	var target Target
	if err := json.NewDecoder(resp.Body).Decode(&target); err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	// Récupérer les repas du jour
	today := time.Now().Format("2006-01-02")
	mealsResp, err := http.Get(fmt.Sprintf("%s/meals/user/%s?date=%s", apiURL, id, today))
	if err != nil {
		fmt.Println("Error getting today's meals:", err)
		return
	}
	defer mealsResp.Body.Close()

	var meals []struct {
		ID    uint   `json:"id"`
		Type  string `json:"type"`
		Foods []Food `json:"foods"`
	}

	if err := json.NewDecoder(mealsResp.Body).Decode(&meals); err != nil {
		fmt.Println("Error parsing meals response:", err)
		return
	}

	// Calculer les totaux des repas
	var consumed struct {
		Calories float64
		Protein  float64
		Carbs    float64
		Fat      float64
		Fiber    float64
	}

	for _, meal := range meals {
		for _, food := range meal.Foods {
			consumed.Calories += food.Calories
			consumed.Protein += food.Protein
			consumed.Carbs += food.Carbs
			consumed.Fat += food.Fat
			consumed.Fiber += food.Fiber
		}
	}

	// Afficher les objectifs et la progression
	fmt.Println("\nDaily Nutrition Targets and Progress:")
	fmt.Printf("Calories: %.0f/%.0f kcal (%.1f%%)\n", consumed.Calories, target.Calories, (consumed.Calories/target.Calories)*100)
	fmt.Printf("Protein:  %.1f/%.1f g (%.1f%%)\n", consumed.Protein, target.Protein, (consumed.Protein/target.Protein)*100)
	fmt.Printf("Carbs:    %.1f/%.1f g (%.1f%%)\n", consumed.Carbs, target.Carbs, (consumed.Carbs/target.Carbs)*100)
	fmt.Printf("Fat:      %.1f/%.1f g (%.1f%%)\n", consumed.Fat, target.Fat, (consumed.Fat/target.Fat)*100)
	fmt.Printf("Fiber:    %.1f/%.1f g (%.1f%%)\n", consumed.Fiber, target.Fiber, (consumed.Fiber/target.Fiber)*100)

	// Afficher les repas du jour
	fmt.Println("\nToday's Meals:")
	for _, meal := range meals {
		fmt.Printf("\n%s:\n", meal.Type)
		for _, food := range meal.Foods {
			fmt.Printf("- %s (%.0f kcal)\n", food.Name, food.Calories)
		}
	}
}

func setTargets(id string) {
	target := promptUserTargets()

	payload, err := json.Marshal(target)
	if err != nil {
		fmt.Println("Error preparing targets:", err)
		return
	}

	resp, err := http.Post(fmt.Sprintf("%s/users/%s/targets", apiURL, id), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error setting targets:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Server returned", resp.Status)
		return
	}

	fmt.Println("Targets set successfully!")
}

func promptUserTargets() Target {
	var target Target

	fmt.Print("Daily Calories Target (kcal): ")
	fmt.Scanf("%f\n", &target.Calories)

	fmt.Print("Daily Protein Target (g): ")
	fmt.Scanf("%f\n", &target.Protein)

	fmt.Print("Daily Carbs Target (g): ")
	fmt.Scanf("%f\n", &target.Carbs)

	fmt.Print("Daily Fat Target (g): ")
	fmt.Scanf("%f\n", &target.Fat)

	fmt.Print("Daily Fiber Target (g): ")
	fmt.Scanf("%f\n", &target.Fiber)

	return target
}

func recordWeight(id string) {
	var weight float64
	var note string

	fmt.Print("Enter your weight (kg): ")
	fmt.Scanf("%f\n", &weight)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a note (optional, press Enter to skip): ")
	note, _ = reader.ReadString('\n')
	note = strings.TrimSpace(note)

	record := WeightRecord{
		Weight: weight,
		Date:   time.Now(),
		Note:   note,
	}

	payload, err := json.Marshal(record)
	if err != nil {
		fmt.Println("Error preparing weight record:", err)
		return
	}

	resp, err := http.Post(fmt.Sprintf("%s/users/%s/weight", apiURL, id), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error recording weight:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Server returned", resp.Status)
		return
	}

	fmt.Println("Weight recorded successfully!")
}

func viewWeightHistory(id string) {
	resp, err := http.Get(fmt.Sprintf("%s/users/%s/weight/history", apiURL, id))
	if err != nil {
		fmt.Println("Error getting weight history:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Server returned", resp.Status)
		return
	}

	var records []WeightRecord
	if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	if len(records) == 0 {
		fmt.Println("No weight records found.")
		return
	}

	fmt.Println("\nWeight History:")
	fmt.Println("Date\t\tWeight\tNote")
	fmt.Println("----------------------------------------")
	for _, record := range records {
		dateStr := record.Date.Format("2006-01-02")
		fmt.Printf("%s\t%.1f kg\t%s\n", dateStr, record.Weight, record.Note)
	}
}
