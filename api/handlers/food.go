package handlers

import (
	"net/http"

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

	c.JSON(http.StatusOK, gin.H{"foods": foods})
}
