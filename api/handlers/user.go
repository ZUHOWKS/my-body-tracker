package handlers

import (
	"net/http"

	"github.com/ZUHOWKS/my-body-tracker/api/models"
	"github.com/ZUHOWKS/my-body-tracker/internal/calculator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUserStats(c *gin.Context) {
	var user models.User
	if err := h.db.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	bmi := calculator.CalculateBMI(user.Weight, user.Height)
	bfp := calculator.CalculateBFP(bmi, user.Age, 1) // Assuming male for now, should be part of user model
	bmr := calculator.CalculateBasalMetabolism(user.Weight, user.Height, user.Age, 1)

	c.JSON(http.StatusOK, gin.H{
		"height": user.Height,
		"weight": user.Weight,
		"bmi":    bmi,
		"bfp":    bfp,
		"bmr":    bmr,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var user models.User
	if err := h.db.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
