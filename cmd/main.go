package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"weather-api/internal/auth"
	"weather-api/internal/client"
	"weather-api/internal/config"
	"weather-api/internal/handler"
	"weather-api/internal/middleware"
	"weather-api/internal/repository"
	"weather-api/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	jwtManager, err := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiration)
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

	weatherService := service.NewWeatherService(
		client.WeatherProviderAdapter{Client: weatherClient},
		client.GeoProviderAdapter{Client: geoClient},
		client.CountryProviderAdapter{Client: countryClient},
	)
	authService := service.NewAuthService(userRepo, jwtManager)
	userService := service.NewUserService(userRepo, cityRepo, historyRepo)
	cityService := service.NewCityService(cityRepo)
	userWeatherService := service.NewUserWeatherService(cityRepo, historyRepo, weatherService)

	if err := authService.EnsureAdminAccount(
		context.Background(),
		cfg.BootstrapAdminEmail,
		cfg.BootstrapAdminPass,
	); err != nil {
		log.Fatal(err)
	}

	weatherHandler := handler.NewWeatherHandler(weatherService)
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	cityHandler := handler.NewCityHandler(cityService)
	userWeatherHandler := handler.NewUserWeatherHandler(userWeatherService)
	authMiddleware := middleware.NewAuthMiddleware(jwtManager, userService)

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

	addr := ":" + cfg.HTTPPort
	log.Printf("server started on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
