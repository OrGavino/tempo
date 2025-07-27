package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country" `
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
	tomorrow := flag.Bool("tomorrow", false, "Get tomorrow's weather forecast")
	flag.Parse()

	days := 1
	if *tomorrow {
		days = 2
	}

	url := fmt.Sprintf(
		"https://api.weatherapi.com/v1/forecast.json?key=f7df7cc30438443f98c162015252207&q=auto:ip&days=%d&aqi=no&alerts=no",
		days,
	)

	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather API not available")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	forecastIndex := 0
	if *tomorrow {
		forecastIndex = 1
	}

	label := "Today"
	if *tomorrow {
		label = "Tomorrow"
	}

	if forecastIndex >= len(weather.Forecast.Forecastday) {
		panic("Forecast for requested day is not available")
	}

	if forecastIndex == 1 {
		fmt.Printf("%s weather will be:\n", label)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[forecastIndex].Hour

	fmt.Printf(
		"🕒 %s, %s: %.0fC, %s\n",
		location.Name,
		location.Country,
		current.TempC,
		current.Condition.Text,
	)

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		if date.Before(time.Now()) {
			continue
		}

		fmt.Printf(
			"%s- %.fC, %.f%%, %s\n",
			date.Format("15:04"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)
	}
}
