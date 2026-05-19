package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"weather-api/internal/auth"
	"weather-api/internal/client"
	"weather-api/internal/config"
	"weather-api/internal/handler"
	"weather-api/internal/middleware"
	"weather-api/internal/repository"
	"weather-api/internal/service"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	jwtManager, err := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiration)
	if err != nil {
		logger.Fatal("failed to create jwt manager", zap.Error(err))
	}

	router := chi.NewRouter()
	router.Use(middleware.NewLoggingMiddleware(logger).Handle)
	router.Use(middleware.NewRecoveryMiddleware(logger).Handle)

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
		logger.Fatal("failed to ensure admin account", zap.Error(err))
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
	logger.Info("server_started", zap.String("addr", addr))

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("server_stopped", zap.Error(err))
	}
}
