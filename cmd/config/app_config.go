package config

import (
	"Go-Starter-Template/internal/api/handlers"
	"Go-Starter-Template/internal/api/routes"
	"Go-Starter-Template/internal/middleware"
	"Go-Starter-Template/internal/utils"
	"Go-Starter-Template/internal/utils/storage"
	"Go-Starter-Template/pkg/jwt"
	"Go-Starter-Template/pkg/midtrans"
	"Go-Starter-Template/pkg/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

func NewApp(db *gorm.DB) (*fiber.App, error) {
	utils.InitValidator()
	app := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
	})
	middlewares := middleware.NewMiddleware()
	jwtService := jwt.NewJWTService()
	validator := utils.Validate

	// setting up logging and limiter
	logDir := "logs"
	logFile := "app.log"
	log_path := filepath.Join(logDir, logFile)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	file, err := os.OpenFile(log_path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	app.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Jakarta",
		Output:     file,
	}))

	// storage
	awsS3 := storage.NewAwsS3()

	// Repository
	userRepository := user.NewUserRepository(db)
	midtransRepository := midtrans.NewMidtransRepository(db)

	// Service
	userService := user.NewUserService(userRepository, awsS3, jwtService)
	midtransService := midtrans.NewMidtransService(
		midtransRepository,
		userRepository,
	)

	// Handler
	userHandler := handlers.NewUserHandler(userService, validator)
	midtransHandler := handlers.NewMidtransHandler(midtransService, validator)

	// routes
	routesConfig := routes.Config{
		App:             app,
		UserHandler:     userHandler,
		MidtransHandler: midtransHandler,
		Middleware:      middlewares,
		JwtService:      jwtService,
	}
	routesConfig.Setup()
	return app, nil
}
