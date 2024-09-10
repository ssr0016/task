package server

import (
	"task/config"
	"task/internal/db"
	"task/internal/logger"

	"task/internal/api/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Server struct {
	app       *fiber.App
	port      string
	db        db.DB
	cfg       *config.Config
	log       *logger.Logger
	jwtSecret string
}

func NewServer(cfg *config.Config) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: errors.DefaultErrorHandler,
	})

	app.Use(cors.New())

	port := ":" + cfg.Port

	// Initialize Database
	sqlxDB := &db.SqlxDB{
		DB: cfg.DB,
	}

	return &Server{
		app:       app,
		port:      port,
		db:        sqlxDB,
		cfg:       cfg,
		log:       cfg.Logger,
		jwtSecret: cfg.JwtSecret,
	}
}

func (s *Server) Start() error {
	s.SetupRoutes()
	return s.app.Listen(s.port)
}

func (s *Server) Stop() error {
	s.log.Sync()
	return s.app.Shutdown()
}
