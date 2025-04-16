package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ZUHOWKS/my-body-tracker/api/handlers"
	"github.com/ZUHOWKS/my-body-tracker/api/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Entrypoint() {
	// Load .env file if it exists
	_ = godotenv.Load() // Ignore error if file doesn't exist

	// Database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	httpClient := &http.Client{}

	db.AutoMigrate(
		&models.User{},
		&models.Food{},
		&models.Meal{},
		&models.Target{},
		&models.WeightRecord{},
	)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(db)
	foodHandler := handlers.NewFoodHandler(db, httpClient)
	mealHandler := handlers.NewMealHandler(db)

	r := gin.Default()

	userRoutes := r.Group("/users")
	{
		userRoutes.GET("/", userHandler.ListUsers)
		userRoutes.POST("/", userHandler.CreateUser)
		userRoutes.GET("/:id", userHandler.GetUser)
		userRoutes.GET("/:id/stats", userHandler.GetUserStats)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.GET("/:id/targets", userHandler.GetUserTargets)
		userRoutes.POST("/:id/targets", userHandler.SetUserTargets)
		userRoutes.POST("/:id/weight", userHandler.RecordWeight)
		userRoutes.GET("/:id/weight/history", userHandler.GetWeightHistory)
	}

	foodRoutes := r.Group("/foods")
	{
		foodRoutes.GET("/search", foodHandler.SearchFood)
	}

	mealRoutes := r.Group("/meals")
	{
		mealRoutes.POST("/", mealHandler.CreateMeal)
		mealRoutes.GET("/user/:userId", mealHandler.GetUserMeals)
		mealRoutes.POST("/:id/foods", mealHandler.AddFoodToMeal)
	}

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
