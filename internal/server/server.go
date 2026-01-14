package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"werewolf-backend/internal/config"
	"werewolf-backend/internal/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
)

type Server struct {
	app *fiber.App
	cfg *config.Config
}

func New(cfg *config.Config) *Server {
	return &Server{
		app: fiber.New(),
		cfg: cfg,
	}
}

func (s *Server) NewServer(cfg *config.Config) *Server {
	app := fiber.New(fiber.Config{
		AppName:               "Werewolf-Game",
		ServerHeader:          "Werewolf-Game",
		DisableStartupMessage: true,
		StrictRouting:         true,
		CaseSensitive:         true,
		ReadTimeout:           cfg.Server.ReadTimeout,
		WriteTimeout:          cfg.Server.WriteTimeout,
		IdleTimeout:           time.Minute,
		ErrorHandler:          customErrorHandler,
	})

	// Security middleware (ป้องกัน common attacks)
	app.Use(helmet.New(helmet.Config{
		XSSProtection:           "1; mode=block",
		ContentTypeNosniff:      "nosniff",
		XFrameOptions:           "DENY",
		ReferrerPolicy:          "no-referrer",
		CrossOriginOpenerPolicy: "same-origin",
	}))

	// Request ID (ติดตาม request แต่ละตัว)
	app.Use(requestid.New(requestid.Config{
		Generator: func() string {
			return uuid.New().String()
		},
	}))

	// Logging
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${locals:requestid} ${status} - ${method} ${path} (${latency})\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Bangkok",
	}))

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     getAllowedOrigins(cfg),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
		MaxAge:           3600,
	}))

	return &Server{
		app: app,
		cfg: cfg,
	}

}

func (s *Server) Run() error {
	// Setup routes
	router.Setup(s.app, s.cfg)

	// Handle graceful shutdown
	return s.runWithGracefulShutdown()
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}
	// Get RequestID safely
	requestID := c.Locals("requestid")
	if requestID == nil {
		requestID = "unknown"
	}
	// Log error
	log.Printf("Error: %v (RequestID: %v)", err, requestID)
	return c.Status(code).JSON(fiber.Map{
		"error":      message,
		"request_id": requestID,
	})
}

func getAllowedOrigins(cfg *config.Config) string {
	if cfg.IsDevelopment() {
		return "http://localhost:3000"
	}
	// Production - เพิ่ม domain จริง
	return "https://your-frontend-domain.com,https://your-app.vercel.app"
}

func (s *Server) runWithGracefulShutdown() error {
	// Start server in goroutine
	addr := fmt.Sprintf("%s:%s", s.cfg.Server.Address, s.cfg.Server.Port)

	go func() {
		log.Printf("🚀 Server starting on %s", addr)
		log.Printf("📊 Environment: %s", s.cfg.Server.Environment)
		if err := s.app.Listen(addr); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()
	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("🛑 Shutting down server...")
	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.Server.ShutdownTimeout)
	defer cancel()
	if err := s.app.ShutdownWithContext(ctx); err != nil {
		log.Printf("❌ Server forced to shutdown: %v", err)
		return err
	}
	log.Println("✅ Server stopped gracefully")

	return nil
}
