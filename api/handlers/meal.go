package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ZUHOWKS/my-body-tracker/api/models"
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
	query := h.db.Where("user_id = ?", c.Param("userId"))

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
		query = query.Where("type = ?", mealType)
	}

	// Execute query with preloaded foods
	if err := query.Preload("Foods").Order("date DESC, type").Find(&meals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, meals)
}

func (h *MealHandler) AddFoodToMeal(c *gin.Context) {
	mealID := c.Param("id")
	var food models.Food
	if err := c.ShouldBindJSON(&food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var meal models.Meal
	if err := h.db.First(&meal, mealID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Meal not found"})
		return
	}

	if err := h.db.Model(&meal).Association("Foods").Append(&food); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, meal)
}
