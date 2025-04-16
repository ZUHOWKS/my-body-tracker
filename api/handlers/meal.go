package handlers

import (
	"net/http"
	"time"

	"github.com/ZUHOWKS/my-body-tracker/api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MealHandler struct {
	db *gorm.DB
}

func NewMealHandler(db *gorm.DB) *MealHandler {
	return &MealHandler{db: db}
}

func (h *MealHandler) CreateMeal(c *gin.Context) {
	var meal models.Meal
	if err := c.ShouldBindJSON(&meal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate meal type
	switch meal.Type {
	case models.Breakfast, models.Lunch, models.Break, models.Dinner:
		// Valid type
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meal type. Must be one of: breakfast, lunch, break, dinner"})
		return
	}

	// Set date to current date if not provided
	if meal.Date.IsZero() {
		meal.Date = time.Now()
	}

	if err := h.db.Create(&meal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, meal)
}

func (h *MealHandler) GetUserMeals(c *gin.Context) {
	var meals []models.Meal
	query := h.db.Model(&models.Meal{}).Where("user_id = ?", c.Param("userId"))

	// Filter by date if provided
	if date := c.Query("date"); date != "" {
		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}
		query = query.Where("DATE(date) = DATE(?)", parsedDate)
	}

	// Filter by type if provided
	if mealType := c.Query("type"); mealType != "" {
		query = query.Where("meal_type = ?", mealType)
	}

	// Execute query with preloaded foods
	if err := query.Preload("Foods").Order("date DESC, meal_type").Find(&meals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, meals)
}

func (h *MealHandler) AddFoodToMeal(c *gin.Context) {
	mealID := c.Param("id")

	// Parse the request body to get the foodId
	var request struct {
		FoodID string `json:"foodId"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.FoodID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "foodId is required"})
		return
	}

	// Vérifier si l'aliment existe déjà dans la base de données
	var existingFood models.Food
	if err := h.db.Where("fdc_id = ?", request.FoodID).First(&existingFood).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Food not found. Please search for it first."})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
			return
		}
	}

	var meal models.Meal
	if err := h.db.Preload("Foods").First(&meal, mealID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meal not found"})
		return
	}

	// Vérifier si l'aliment est déjà dans le repas
	var foodExists bool
	for _, food := range meal.Foods {
		if food.FdcID == request.FoodID {
			foodExists = true
			break
		}
	}

	if foodExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Food is already in this meal"})
		return
	}

	// Ajouter l'aliment au repas
	if err := h.db.Model(&meal).Association("Foods").Append(&existingFood); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add food to meal: " + err.Error()})
		return
	}

	// Recharger le repas avec ses aliments
	if err := h.db.Preload("Foods").First(&meal, mealID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload meal: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, meal)
}
