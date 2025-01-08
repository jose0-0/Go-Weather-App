package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type GeoResponse struct {
	Results []LatLong `json:"results"`
}

type LatLong struct {
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func main() {
	latLong, err := getLatLong("London")
	if err != nil {
		log.Fatalf("Failed to get latitude and longitude: %s", err)
	}
	fmt.Printf("Latitude: %f, Longitude: %f\n", latLong.Latitude, latLong.Longitude)

	weather, err := getWeather(*latLong)
	if err != nil {
		log.Fatalf("Failed to get weather: %s", err)
	}
	fmt.Printf("Weather: %s\n", weather)
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