package main

import (
	"log"

	"github.com/Naitik-ag/fleetintel-backend/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "FleetIntel API v1.0",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path}\n",
	}))

	api := app.Group("/api")

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "FleetIntel API is running",
		})
	})

	trackingHandler := handlers.NewTrackingHandler()

	tracking := api.Group("/tracking")
	tracking.Post("/location", trackingHandler.UpdateLocation)
	tracking.Get("/vehicles/:id/live", trackingHandler.GetLiveLocation)
	tracking.Get("/vehicles/:id/history", trackingHandler.GetLocationHistory)

	log.Println("Server starting on http://localhost:8001")
	err := app.Listen(":8001")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
