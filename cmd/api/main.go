package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"CLOAKBE/internal/config"
	"CLOAKBE/internal/database"
	"CLOAKBE/internal/handler"
	"CLOAKBE/internal/middleware"
	"CLOAKBE/internal/repository"
	"CLOAKBE/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// @title           CLOAK API
// @version         1.0
// @description     Digital Ticketing System API
// @termsOfService  http://swagger.io/terms/

// @contact.name   Support
// @contact.email  support@cloak.local

// @license.name  MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Init repositories
	businessRepo := repository.NewPostgresBusinessRepository(db)
	customerRepo := repository.NewPostgresCustomerRepository(db)
	serviceRepo := repository.NewPostgresServiceRepository(db)
	slotRepo := repository.NewPostgresSlotRepository(db)
	ticketRepo := repository.NewPostgresTicketRepository(db)

	// Init usecases
	authUsecase := usecase.NewAuthUsecase(businessRepo, customerRepo, cfg.JWTSecret)
	ticketUsecase := usecase.NewTicketUsecase(ticketRepo, slotRepo, serviceRepo, businessRepo)
	serviceUsecase := usecase.NewServiceUsecase(serviceRepo, slotRepo, businessRepo)

	// Init handlers
	authHandler := handler.NewAuthHandler(authUsecase)
	ticketHandler := handler.NewTicketHandler(ticketUsecase)
	serviceHandler := handler.NewServiceHandler(serviceUsecase)

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "CLOAK API v1.0",
		BodyLimit:    1024 * 1024,
		ErrorHandler: defaultErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin,Content-Type,Authorization,Accept",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Public routes
	public := app.Group("/api/v1")
	public.Post("/auth/business/register", authHandler.BusinessRegister)
	public.Post("/auth/business/login", authHandler.BusinessLogin)
	public.Post("/auth/customer/login", authHandler.CustomerLogin)

	// Protected routes (require JWT)
	protected := app.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// Business routes (role: business)
	business := protected.Group("/tickets")
	business.Use(middleware.RoleMiddleware("business"))
	business.Post("/checkin", ticketHandler.CheckIn)
	business.Post("/scan", ticketHandler.Scan)
	business.Post("/:id/release", ticketHandler.Release)

	// Service routes (role: business)
	services := protected.Group("/services")
	services.Use(middleware.RoleMiddleware("business"))
	services.Post("", serviceHandler.CreateService)
	services.Get("", serviceHandler.ListServices)
	services.Get("/:id", serviceHandler.GetService)
	services.Get("/:id/stats", serviceHandler.GetServiceStats)

	// Customer routes (role: customer)
	customer := protected.Group("/tickets")
	customer.Use(middleware.RoleMiddleware("customer"))
	customer.Get("/:id", ticketHandler.GetTicket)

	// Start server with graceful shutdown
	go func() {
		log.Printf("Starting server on port %s", cfg.ServerPort)
		if err := app.Listen(":" + cfg.ServerPort); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func defaultErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "Internal Server Error"

	if fe, ok := err.(*fiber.Error); ok {
		code = fe.Code
		msg = fe.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"code":    code,
		"message": msg,
		"error":   err.Error(),
	})
}
