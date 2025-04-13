package handlers

import (
	"net/http"

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

	if err := h.db.Create(&meal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, meal)
}

func (h *MealHandler) GetUserMeals(c *gin.Context) {
	var meals []models.Meal
	if err := h.db.Where("user_id = ?", c.Param("userId")).Find(&meals).Error; err != nil {
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
