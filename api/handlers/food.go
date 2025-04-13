package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ZUHOWKS/my-body-tracker/api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FoodHandler struct {
	db *gorm.DB
}

type FDCResponse struct {
	Foods []struct {
		FdcId         string  `json:"fdcId"`
		Description   string  `json:"description"`
		ServingSize   float64 `json:"servingSize"`
		FoodNutrients []struct {
			NutrientName string  `json:"nutrientName"`
			Value        float64 `json:"value"`
		} `json:"foodNutrients"`
	} `json:"foods"`
}

func NewFoodHandler(db *gorm.DB) *FoodHandler {
	return &FoodHandler{db: db}
}

func (h *FoodHandler) SearchFood(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	apiKey := os.Getenv("FDC_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "FDC API key not configured"})
		return
	}

	url := fmt.Sprintf("https://api.nal.usda.gov/fdc/v1/foods/search?api_key=%s&query=%s", apiKey, query)
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch from FDC API"})
		return
	}
	defer resp.Body.Close()

	var fdcResp FDCResponse
	if err := json.NewDecoder(resp.Body).Decode(&fdcResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse FDC response"})
		return
	}

	c.JSON(http.StatusOK, fdcResp)
}

func (h *FoodHandler) SaveFood(c *gin.Context) {
	var food models.Food
	if err := c.ShouldBindJSON(&food); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&food).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, food)
}
