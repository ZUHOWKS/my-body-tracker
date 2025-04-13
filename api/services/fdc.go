package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

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

func SearchFood(query string, client http.Client) (FDCResponse, error) {
	apiKey := os.Getenv("FDC_API_KEY")
	if apiKey == "" {
		return FDCResponse{}, fmt.Errorf("FDC API key not configured")
	}

	baseURL := "https://api.nal.usda.gov/fdc/v1/foods/search"
	params := url.Values{}
	params.Add("api_key", apiKey)
	params.Add("query", query)

	fullURL := baseURL + "?" + params.Encode()

	resp, err := client.Get(fullURL)
	if err != nil {
		return FDCResponse{}, fmt.Errorf("failed to fetch FDC data: %s", err.Error())
	}
	defer resp.Body.Close()

	var fdcResp FDCResponse
	if err := json.NewDecoder(resp.Body).Decode(&fdcResp); err != nil {
		return FDCResponse{}, fmt.Errorf("failed to parse FDC response: %w", err)
	}

	return fdcResp, nil
}
