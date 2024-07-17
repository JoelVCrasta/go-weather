package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Weather struct
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
		Humidity   int     `json:"humidity"`
		FeelslikeC float64 `json:"feelslike_c"`
	} `json:"current"`

	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				Humidity     int     `json:"humidity"`
				FeelslikeC   float64 `json:"feelslike_c"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=b78509e46e9243ffbc953804241707&q=Mangalore&days=1&aqi=no&alerts=no")

	if err != nil {
		panic(err)
	}
	defer res.Body.Close() // Close the response body

	if res.StatusCode != 200 {
		panic("Weather api not available")
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	// Unmarshal the response body to Weather struct
	var weather Weather
	err = json.Unmarshal(body, &weather)

	if err != nil {
		panic(err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	fmt.Println("<------- Current Weather ------->")
	fmt.Printf("%-12s %-12s %-8s %-12s %-8s %s\n", "City", "Country", "Temp", "FeelsLike", "Humidity", "Condition")
	fmt.Printf("%-12s %-12s %-8s %-12s %-8v %s\n",
		location.Name,
		location.Country,
		fmt.Sprintf("%.1fC", current.TempC),
		fmt.Sprintf("%.1fC", current.FeelslikeC),
		current.Humidity,
		current.Condition.Text,
	)

	fmt.Println("\n<------- Forecast  ------->")
	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		if date.Before(time.Now()) {
			continue
		}

		fmt.Printf(
			"%s - %.0fC, %.0f%%, %s\n",
			date.Format("15:04"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)
	}
}
