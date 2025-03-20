package config

import (
	"Go-Starter-Template/internal/api/handlers"
	"Go-Starter-Template/internal/api/routes"
	"Go-Starter-Template/internal/middleware"
	"Go-Starter-Template/internal/utils"
	"Go-Starter-Template/internal/utils/storage"
	"Go-Starter-Template/pkg/chat"
	"Go-Starter-Template/pkg/company"
	"Go-Starter-Template/pkg/job"
	"Go-Starter-Template/pkg/jwt"
	"Go-Starter-Template/pkg/midtrans"
	"Go-Starter-Template/pkg/notification"
	"Go-Starter-Template/pkg/post"
	"Go-Starter-Template/pkg/user"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
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
	companyRepository := company.NewCompanyRepository(db)
	midtransRepository := midtrans.NewMidtransRepository(db)
	jobRepository := job.NewJobRepository(db)
	chatRepository := chat.NewChatRepository(db)
	notificationRepository := notification.NewNotificationRepository(db)
	postRepository := post.NewPostRepository(db)

	// Service
	userService := user.NewUserService(userRepository, awsS3, jwtService)
	companyService := company.NewCompanyService(companyRepository, awsS3, jwtService)
	midtransService := midtrans.NewMidtransService(
		midtransRepository,
		userRepository,
	)
	jobService := job.NewJobService(jobRepository, awsS3, jwtService)
	chatService := chat.NewChatService(chatRepository, jwtService)
	notificationService := notification.NewNotificationService(notificationRepository, jwtService)
	postService := post.NewPostService(postRepository, awsS3, jwtService)

	// Handler
	userHandler := handlers.NewUserHandler(userService, validator)
	companyHandler := handlers.NewCompanyHandler(companyService, validator)
	midtransHandler := handlers.NewMidtransHandler(midtransService, validator)
	jobHandler := handlers.NewJobHandler(jobService, validator)
	chatServerHandler := handlers.NewChatServerHandler()
	chatHandler := handlers.NewChatHandler(chatService, validator)
	notificationHandler := handlers.NewNotificationHandler(notificationService, validator)
	postHandler := handlers.NewPostHandler(postService, validator)

	// routes
	routesConfig := routes.Config{
		App:                 app,
		UserHandler:         userHandler,
		CompanyHandler:      companyHandler,
		MidtransHandler:     midtransHandler,
		Middleware:          middlewares,
		JwtService:          jwtService,
		JobHandler:          jobHandler,
		ChatServerHandler:   *chatServerHandler,
		ChatHandler:         chatHandler,
		NotificationHandler: notificationHandler,
		PostHandler:         postHandler,
	}

	routesConfig.Setup()
	return app, nil
}
