package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	// "log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type GeoResponse struct {
	Results []LatLong `json:"results"`
}

type LatLong struct {
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func main() {
    r := gin.Default()

    r.GET("/weather", func(c *gin.Context) {
   	 city := c.Query("city")
   	 latlong, err := getLatLong(city)
   	 if err != nil {
   		 c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
   		 return
   	 }

   	 weather, err := getWeather(*latlong)
   	 if err != nil {
   		 c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
   		 return
   	 }

   	 c.JSON(http.StatusOK, gin.H{"weather": weather})
    })

    r.Run()
}

func getLatLong(city string) (*LatLong, error) {
	endpoint := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=en&format=json", url.QueryEscape(city))
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("Error making request to Geo API: %w", err)
	}
	defer resp.Body.Close()

	var response GeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("Error Decoding response: %w", err)
	}

	if len(response.Results) < 1 {
		return nil, errors.New("No results found")
	}

	return &response.Results[0], nil
}

func getWeather(latLong LatLong) (string, error) {
	endpoint := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.6f&longitude=%.6f&hourly=temperature_2m", latLong.Latitude, latLong.Longitude)
	resp, err := http.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("Error making request to Weather API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body: %w", err)
	}

	return string(body), nil

}