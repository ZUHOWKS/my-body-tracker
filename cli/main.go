package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const apiURL = "http://localhost:8080"

func Entrypoint() {
	fmt.Println("Welcome to My Body Tracker CLI!")
	fmt.Println("Available commands:")
	fmt.Println("  profile create - Create a new user profile")
	fmt.Println("  profile view <id> - View user profile and stats")
	fmt.Println("  profile list - List available users")
	fmt.Println("  profile select <id> - Select a user")
	fmt.Println("  food search <query> - Search for food in database")
	fmt.Println("  meal add <type> <date> <food_name> - Add food to meal")
	fmt.Println("  meal view <type> <date> - View meal details and nutrients")
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
		handleFoodCommand(args)

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
