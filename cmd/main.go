package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"weather-api/internal/auth"
	"weather-api/internal/client"
	"weather-api/internal/handler"
	"weather-api/internal/middleware"
	"weather-api/internal/repository"
	"weather-api/internal/service"
)

func main() {
	jwtSecret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	jwtManager, err := auth.NewJWTManager(jwtSecret, 24*time.Hour)
	if err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	weatherClient := client.NewWeatherClient(httpClient)
	geoClient := client.NewGeoClient(httpClient)
	countryClient := client.NewCountryClient(httpClient)

	userRepo := repository.NewInMemoryUserRepository()
	cityRepo := repository.NewInMemoryCityRepository()
	historyRepo := repository.NewInMemoryWeatherHistoryRepository()

	weatherService := service.NewWeatherService(weatherClient, geoClient, countryClient)
	authService := service.NewAuthService(userRepo, jwtManager, parseCSVEnv("ADMIN_EMAILS"))
	userService := service.NewUserService(userRepo, cityRepo, historyRepo)
	cityService := service.NewCityService(cityRepo)
	userWeatherService := service.NewUserWeatherService(cityRepo, historyRepo, weatherService)

	weatherHandler := handler.NewWeatherHandler(weatherService)
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	cityHandler := handler.NewCityHandler(cityService)
	userWeatherHandler := handler.NewUserWeatherHandler(userWeatherService)
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	router.Route("/weather", func(r chi.Router) {
		r.With(authMiddleware.Handle).Get("/", userWeatherHandler.GetCurrent)
		r.With(authMiddleware.Handle).Get("/history", userWeatherHandler.GetHistory)
		r.Get("/{city}", weatherHandler.GetWeatherByCity)
		r.Get("/country/{country}", weatherHandler.GetWeatherByCountry)
		r.Get("/country/{country}/top", weatherHandler.GetTopCitiesByCountry)
	})

	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	router.Group(func(r chi.Router) {
		r.Use(authMiddleware.Handle)

		r.Route("/cities", func(r chi.Router) {
			r.Post("/", cityHandler.Create)
			r.Get("/", cityHandler.List)
			r.Delete("/{city_id}", cityHandler.Delete)
		})

		r.Route("/users", func(r chi.Router) {
			r.Get("/me", userHandler.Me)

			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole("admin"))
				r.Get("/", userHandler.List)
				r.Get("/{id}", userHandler.GetByID)
				r.Delete("/{id}", userHandler.Delete)
			})
		})
	})

	addr := ":8080"
	log.Printf("server started on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}

func parseCSVEnv(name string) []string {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if value := strings.TrimSpace(part); value != "" {
			result = append(result, value)
		}
	}

	return result
}
