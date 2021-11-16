package routes

import (
	"context"
	"io/ioutil"
	"os"

	swagger "github.com/arsmn/fiber-swagger/v2"

	"github.com/casbin/casbin/v2"
	"github.com/fusesuphasin/go-fiber/app/controller"
	"github.com/fusesuphasin/go-fiber/app/interfaces"
	"github.com/fusesuphasin/go-fiber/app/middleware"
	"github.com/fusesuphasin/go-fiber/app/repository"
	"github.com/fusesuphasin/go-fiber/app/service"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoute(app *fiber.App, ctx context.Context, log interfaces.Logger, enforcer *casbin.Enforcer) {
	UserController := controller.NewUserController(log)
	//HomeController := controller.NewHomeController(log, app)
	RoleController := controller.NewRoleController(log, app, enforcer)
	PermissionController := controller.NewPermissionController(log, enforcer, app)
	QueueController := controller.NewQueueController(log)
	UploadController := controller.NewUploadController(log, app, enforcer)
	app.Get("", func(c *fiber.Ctx) error {
		return c.JSON("Wellcome 2")
	})

	app.Get("/docs/*", swagger.New(
		swagger.Config{
			URL:         os.Getenv("APP_URL") + "/swagger.json",
			DeepLinking: true,
			// Expand ("list") or Collapse ("none") tag groups by default
			DocExpansion: "none",
		},
	)) // default

	app.Get("/swagger.json", func(c *fiber.Ctx) error {
		body, err := ioutil.ReadFile("docs/swagger.json")
		if err != nil {
			panic(err)
		}
		return c.SendString(string(body))
	})

	//app.Get("/dashboard", monitor.New())

	app.Post("/register", UserController.Register)
	app.Post("/login", UserController.Login)


	
	//get Token to check permission
	app.Use(middleware.JWTProtected(service.UserService{
		UserRepository: repository.UserRepository{
			Ctx: ctx,
		},
	}))
	
	PermissionController.PermissionRouter()
	UploadController.UploadRouter()
	RoleController.RoleRouter()

	//determind role to accessing 
	enforcer.AddPolicy("admin", "role", "manage")
	enforcer.AddPolicy("admin", "permission", "manage")
	enforcer.AddPolicy("admin", "upload", "manage")

	//HomeController.HomeRouter()
	

	app.Get("/test", QueueController.TestQueue)
	app.Get("/testgetQ", QueueController.TestGetFromQueue)

}  