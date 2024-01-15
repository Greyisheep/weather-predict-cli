package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

// Weather struct represents the structure of the JSON response from the weather API.
type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	// Default location is set to Lagos
	q := "Lagos"

	// Check if a location is provided as a command-line argument
	if len(os.Args) >= 2 {
		q = os.Args[1]
	}

	// Make a request to the weather API
	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=d481907ed4154e4fbd1211232241401&q=" + q + "&days=1&aqi=no&alerts=no")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	// Check if the API responds successfully (status code 200)
	if res.StatusCode != 200 {
		panic("Weather API not responding")
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// Unmarshal the JSON response into the Weather struct
	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	// Print the current weather information
	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour
	fmt.Printf("%s, %s: %.0fC, %s\n", location.Name, location.Country, current.TempC, current.Condition.Text)

	// Iterate over the forecast hours and display relevant information
	for _, hour := range hours {
		// Convert the Unix timestamp to a readable date
		date := time.Unix(hour.TimeEpoch, 0)

		// Skip past hours that have already occurred
		if date.Before(time.Now()) {
			continue
		}

		// Format and print the hourly weather information
		message := fmt.Sprintf("%s - %.0fC, %.0f%%, %s\n", date.Format("15:04"), hour.TempC, hour.ChanceOfRain, hour.Condition.Text)

		// Use color package to highlight high chance of rain in red
		if hour.ChanceOfRain < 40 {
			fmt.Print(message)
		} else {
			color.Red(message)
		}
	}
}
