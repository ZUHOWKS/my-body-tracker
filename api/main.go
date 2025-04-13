package api

import (
	"log"
	"net/http"
	"os"

	"github.com/ZUHOWKS/my-body-tracker/api/handlers"
	"github.com/ZUHOWKS/my-body-tracker/api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Entrypoint() {
	// Database connection
	dsn := os.ExpandEnv("host=postgres user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=disable")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	httpClient := &http.Client{}

	db.AutoMigrate(
		&models.User{},
		&models.MealPlan{},
		&models.Meal{},
		&models.Food{},
		&models.DailyIntake{},
	)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(db)
	foodHandler := handlers.NewFoodHandler(db, httpClient)
	mealHandler := handlers.NewMealHandler(db)

	r := gin.Default()

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/", userHandler.CreateUser)
		userRoutes.GET("/:id", userHandler.GetUserStats)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
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
