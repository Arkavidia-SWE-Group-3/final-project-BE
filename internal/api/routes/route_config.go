package routes

import (
	"Go-Starter-Template/internal/api/handlers"
	"Go-Starter-Template/internal/middleware"
	jwtService "Go-Starter-Template/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

type Config struct {
	App               *fiber.App
	UserHandler       handlers.UserHandler
	CompanyHandler    handlers.CompanyHandler
	ChatServerHandler handlers.ChatServerHandler
	ChatHandler       handlers.ChatHandler
	JobHandler        handlers.JobHandler
	MidtransHandler   handlers.MidtransHandler
	Middleware        middleware.Middleware
	JwtService        jwtService.JWTService
}

func (c *Config) Setup() {
	c.App.Use(c.Middleware.CORSMiddleware())
	c.User()
	c.Company()
	c.Job()
	c.Chat()
	c.GuestRoute()
	c.AuthRoute()
}

func (c *Config) User() {
	user := c.App.Group("/api/user")
	{
		user.Get("/search", c.UserHandler.SearchUser)
		user.Post("/register", c.UserHandler.RegisterUser)
		user.Post("/login", c.UserHandler.Login)
		user.Get("/profile/:slug", c.UserHandler.GetProfile)
		user.Post("/update-profile", c.Middleware.AuthMiddleware(c.JwtService), c.Middleware.OnlyAllow("user"), c.UserHandler.UpdateProfile)

		education := user.Group("/education")
		{
			education.Post("/add-education", c.Middleware.AuthMiddleware(c.JwtService), c.Middleware.OnlyAllow("user"), c.UserHandler.PostEducation)
			education.Patch("/update-education", c.Middleware.AuthMiddleware(c.JwtService), c.Middleware.OnlyAllow("user"), c.UserHandler.UpdateEducation)
			education.Delete("/delete-education/:id", c.Middleware.AuthMiddleware(c.JwtService), c.Middleware.OnlyAllow("user"), c.UserHandler.DeleteEducation)
		}

		experience := user.Group("/experience")
		{
			experience.Patch("/update-experience", c.Middleware.AuthMiddleware(c.JwtService), c.Middleware.OnlyAllow("user"), c.UserHandler.UpdateExperience)
			experience.Post("/add-experience", c.Middleware.AuthMiddleware(c.JwtService), c.Middleware.OnlyAllow("user"), c.UserHandler.PostExperience)
			experience.Delete("/delete-experience/:id", c.Middleware.AuthMiddleware(c.JwtService), c.Middleware.OnlyAllow("user"), c.UserHandler.DeleteExperience)
		}

		skills := user.Group("/skills")
		{
			skills.Post("/add-skill", c.Middleware.AuthMiddleware(c.JwtService), c.Middleware.OnlyAllow("user"), c.UserHandler.PostSkill)
			skills.Delete("/delete-skill/:id", c.Middleware.AuthMiddleware(c.JwtService), c.Middleware.OnlyAllow("user"), c.UserHandler.DeleteSkill)
		}

		user.Post("/subscribe", c.Middleware.AuthMiddleware(c.JwtService), c.MidtransHandler.CreateTransaction)

	}

}

func (c *Config) Company() {
	company := c.App.Group("/api/company")
	{
		company.Post("/login", c.CompanyHandler.LoginCompany)
		company.Post("/register", c.CompanyHandler.RegisterCompany)
		company.Get("/profile/:slug", c.CompanyHandler.GetProfile)
		company.Get("/list", c.CompanyHandler.GetListCompany)
		company.Patch("/update-profile", c.Middleware.AuthMiddleware(c.JwtService), c.Middleware.OnlyAllow("company"), c.CompanyHandler.UpdateProfile)
		company.Post("/add-job", c.Middleware.AuthMiddleware(c.JwtService), c.Middleware.OnlyAllow("company"), c.CompanyHandler.AddJob)
		company.Patch("/update-job", c.Middleware.AuthMiddleware(c.JwtService), c.Middleware.OnlyAllow("company"), c.CompanyHandler.UpdateJob)
	}
}

func (c *Config) Job() {
	job := c.App.Group("/api/job")
	{
		job.Get("/detail/:id", c.JobHandler.GetJobDetail)
		job.Get("/search", c.JobHandler.SearchJob)
		job.Get("/applicants/:id", c.Middleware.AuthMiddleware(c.JwtService), c.Middleware.OnlyAllow("company"), c.JobHandler.GetApplicants)
		job.Post("/apply", c.Middleware.AuthMiddleware(c.JwtService), c.Middleware.OnlyAllow("user"), c.JobHandler.ApplyJob)
		job.Post("/update-application", c.Middleware.AuthMiddleware(c.JwtService), c.Middleware.OnlyAllow("company"), c.JobHandler.ChangeApplicationStatus)
	}
}

func (c *Config) Chat() {
	c.ChatServerHandler.SetupRoutes(c.App)

	chat := c.App.Group("/api/chat")
	{
		chat.Get("/rooms", c.Middleware.AuthMiddleware(c.JwtService), c.ChatHandler.GetChatRooms)
		chat.Get("/room/:id", c.Middleware.AuthMiddleware(c.JwtService), c.ChatHandler.GetChatRoom)
		chat.Post("/send", c.Middleware.AuthMiddleware(c.JwtService), c.ChatHandler.SendMessage)
		chat.Get("/messages/:id", c.Middleware.AuthMiddleware(c.JwtService), c.ChatHandler.GetMessages)
	}

}

func (c *Config) GuestRoute() {
	c.App.Get("/api/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "pong, its works. please"})
	})
	c.App.Post("/webhook/midtrans", c.MidtransHandler.MidtransWebhookHandler)
}

func (c *Config) AuthRoute() {
	c.App.Get("/restricted", c.Middleware.AuthMiddleware(c.JwtService), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Access granted"})
	})
	c.App.Get("/me", c.Middleware.AuthMiddleware(c.JwtService), func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		role := c.Locals("role")
		return c.JSON(fiber.Map{
			"message": "Welcome to your dashboard",
			"user_id": userID,
			"role":    role,
		})
	})
}
