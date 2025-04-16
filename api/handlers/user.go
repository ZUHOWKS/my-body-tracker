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

func (h *UserHandler) ListUsers(c *gin.Context) {
	var users []models.User
	if err := h.db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUserStats(c *gin.Context) {
	var user models.User
	if err := h.db.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	bmi := calculator.CalculateBMI(user.Weight, user.Height)
	bfp := calculator.CalculateBFP(bmi, user.Age, user.Sex)
	img := calculator.CalculateIMG(user.Weight, user.Height, user.Sex)
	bmr := calculator.CalculateBasalMetabolism(user.Weight, user.Height, user.Age, user.Sex)

	// Log the response to debug
	response := gin.H{
		"height": user.Height,
		"weight": user.Weight,
		"bmi":    bmi,
		"bfp":    bfp,
		"img":    img,
		"bmr":    bmr,
	}

	// Print to server logs
	c.JSON(http.StatusOK, response)
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

func (h *UserHandler) GetUser(c *gin.Context) {
	var user models.User
	if err := h.db.First(&user, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) SetUserTargets(c *gin.Context) {
	userID := c.Param("id")

	var target models.Target
	if err := c.ShouldBindJSON(&target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Vérifier si l'utilisateur existe
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Définir l'ID de l'utilisateur
	target.UserID = user.ID

	// Vérifier si des objectifs existent déjà pour cet utilisateur
	var existingTarget models.Target
	result := h.db.Where("user_id = ?", user.ID).First(&existingTarget)

	if result.Error == gorm.ErrRecordNotFound {
		// Créer de nouveaux objectifs
		if err := h.db.Create(&target).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else if result.Error != nil {
		// Une erreur s'est produite
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	} else {
		// Mettre à jour les objectifs existants
		target.ID = existingTarget.ID
		if err := h.db.Save(&target).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, target)
}

func (h *UserHandler) GetUserTargets(c *gin.Context) {
	userID := c.Param("id")

	var target models.Target
	if err := h.db.Where("user_id = ?", userID).First(&target).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "No targets found for this user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, target)
}
