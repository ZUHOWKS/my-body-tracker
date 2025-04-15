package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ZUHOWKS/my-body-tracker/api/models"
)

type FDCNutrient struct {
	NutrientName string  `json:"nutrientName"`
	Value        float64 `json:"value"`
}

type FDCFood struct {
	FdcId         int           `json:"fdcId"`
	Description   string        `json:"description"`
	ServingSize   float64       `json:"servingSize"`
	FoodNutrients []FDCNutrient `json:"foodNutrients"`
}

type FDCResponse struct {
	Foods []FDCFood `json:"foods"`
}

// SearchFood searches FDC database for a given query and returns our Food model
func SearchFood(query string, client http.Client) ([]models.Food, error) {
	baseURL := "https://api.nal.usda.gov/fdc/v1/foods/search"
	params := url.Values{}
	params.Add("api_key", "lVglma2Dy1h69QzmRovFef2yOxqABWT0bldH8iLm")
	params.Add("query", query)
	params.Add("pageSize", "5") // Limit results to 5 items

	fullURL := baseURL + "?" + params.Encode()

	resp, err := client.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch FDC data: %s", err.Error())
	}
	defer resp.Body.Close()

	var fdcResp FDCResponse
	if err := json.NewDecoder(resp.Body).Decode(&fdcResp); err != nil {
		return nil, fmt.Errorf("failed to parse FDC response: %w", err)
	}

	// Convert FDC foods to our Food models
	foods := make([]models.Food, len(fdcResp.Foods))
	for i, fdcFood := range fdcResp.Foods {
		protein, carbs, fat, calories, fiber := 0.0, 0.0, 0.0, 0.0, 0.0
		for _, nutrient := range fdcFood.FoodNutrients {
			switch nutrient.NutrientName {
			case "Protein":
				protein = nutrient.Value
			case "Carbohydrate, by difference":
				carbs = nutrient.Value
			case "Total lipid (fat)":
				fat = nutrient.Value
			case "Energy":
				calories = nutrient.Value
			case "Fiber, total dietary":
				fiber = nutrient.Value
			}
		}
		foods[i] = models.Food{
			FdcID:       fmt.Sprintf("%d", fdcFood.FdcId), // Convert int to string
			Name:        fdcFood.Description,
			Protein:     protein,
			Carbs:       carbs,
			Fat:         fat,
			Calories:    calories,
			Fiber:       fiber,
			ServingSize: fdcFood.ServingSize,
		}
	}

	return foods, nil
}
