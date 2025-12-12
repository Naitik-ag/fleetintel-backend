package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"github.com/Naitik-ag/fleetintel-backend/internal/database"
	"github.com/Naitik-ag/fleetintel-backend/internal/handlers"
	"github.com/Naitik-ag/fleetintel-backend/internal/models"
	"github.com/Naitik-ag/fleetintel-backend/internal/repository"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	err = database.ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = database.AutoMigrate(
		&models.Location{},
		&models.DeliveryAgent{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	locationRepo := repository.NewLocationRepository(database.GetDB())
	agentRepo := repository.NewAgentRepository(database.GetDB())

	trackingHandler := handlers.NewTrackingHandler(locationRepo)
	agentHandler := handlers.NewAgentHandler(agentRepo)

	app := fiber.New(fiber.Config{
		AppName: "FleetIntel API",
	})

	app.Use(logger.New())
	app.Use(cors.New())

	setupRoutes(app, trackingHandler, agentHandler)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8001"
	}

	log.Printf("Server is running on http://localhost:%s\n", port)
	log.Fatal(app.Listen(":" + port))
}

func setupRoutes(app *fiber.App, trackingHandler *handlers.TrackingHandler, agentHandler *handlers.AgentHandler) {

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":   "healthy",
			"service":  "FleetIntel API",
			"database": "connected",
		})
	})

	api := app.Group("/api")

	agents := api.Group("/agents")
	agents.Post("/", agentHandler.RegisterAgent)
	agents.Get("/", agentHandler.ListAgents)
	agents.Get("/:id", agentHandler.GetAgent)
	agents.Put("/:id", agentHandler.UpdateAgent)
	agents.Delete("/:id", agentHandler.DeleteAgent)
	agents.Patch("/:id/status", agentHandler.UpdateAgentStatus)
	agents.Get("/:id/stats", agentHandler.GetAgentStats)

	tracking := api.Group("/tracking")
	tracking.Post("/location", trackingHandler.UpdateLocation)
	tracking.Get("/location/:id", trackingHandler.GetLiveLocation)
	tracking.Get("/history/:id", trackingHandler.GetLocationHistory)
}