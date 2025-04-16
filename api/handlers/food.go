package handlers

import (
	"net/http"

	"github.com/ZUHOWKS/my-body-tracker/api/models"
	"github.com/ZUHOWKS/my-body-tracker/api/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FoodHandler struct {
	db         *gorm.DB
	httpClient *http.Client
}

func NewFoodHandler(db *gorm.DB, client *http.Client) *FoodHandler {
	return &FoodHandler{db: db, httpClient: client}
}

func (h *FoodHandler) SearchFood(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	foods, err := services.SearchFood(query, *h.httpClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Enregistrer les aliments dans la base de données locale
	for i := range foods {
		// Vérifier si l'aliment existe déjà
		var existingFood models.Food
		result := h.db.Where("fdc_id = ?", foods[i].FdcID).First(&existingFood)

		if result.Error == gorm.ErrRecordNotFound {
			// L'aliment n'existe pas, le créer
			if err := h.db.Create(&foods[i]).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save food: " + err.Error()})
				return
			}
		} else if result.Error != nil {
			// Une erreur s'est produite lors de la recherche
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + result.Error.Error()})
			return
		}
		// Si l'aliment existe déjà, on utilise celui de la base de données
	}

	c.JSON(http.StatusOK, gin.H{"foods": foods})
}
