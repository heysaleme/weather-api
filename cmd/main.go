package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"weather-api/internal/client"
	"weather-api/internal/handler"
	"weather-api/internal/service"
)

func main() {
	router := chi.NewRouter()

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	weatherClient := client.NewWeatherClient(httpClient)
	geoClient := client.NewGeoClient(httpClient)
	countryClient := client.NewCountryClient(httpClient)

	weatherService := service.NewWeatherService(weatherClient, geoClient, countryClient)
	weatherHandler := handler.NewWeatherHandler(weatherService)

	router.Route("/weather", func(r chi.Router) {
		r.Get("/{city}", weatherHandler.GetWeatherByCity)
		r.Get("/country/{country}", weatherHandler.GetWeatherByCountry)
		r.Get("/country/{country}/top", weatherHandler.GetTopCitiesByCountry)
	})

	addr := ":8080"
	log.Printf("server started on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
