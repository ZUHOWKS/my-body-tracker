package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ZUHOWKS/my-body-tracker/api/handlers"
	"github.com/ZUHOWKS/my-body-tracker/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Database connection
	dsn := os.ExpandEnv("host=postgres user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=disable")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	db.AutoMigrate(
		&models.User{},
		&models.MealPlan{},
		&models.Meal{},
		&models.Food{},
		&models.DailyIntake{},
	)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(db)
	foodHandler := handlers.NewFoodHandler(db)
	mealHandler := handlers.NewMealHandler(db)

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// User routes
	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/", userHandler.CreateUser)
		userRoutes.GET("/:id", userHandler.GetUserStats)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
	}

	// Food routes
	foodRoutes := r.Group("/foods")
	{
		foodRoutes.GET("/search", foodHandler.SearchFood)
		foodRoutes.POST("/", foodHandler.SaveFood)
	}

	// Meal routes
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
