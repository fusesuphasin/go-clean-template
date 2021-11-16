package controller

import (
	"context"

	"github.com/casbin/casbin/v2"
	"github.com/fusesuphasin/go-fiber/app/interfaces"
	"github.com/fusesuphasin/go-fiber/app/middleware"
	"github.com/fusesuphasin/go-fiber/app/repository"
	"github.com/fusesuphasin/go-fiber/app/rules"
	"github.com/fusesuphasin/go-fiber/app/service"
	"github.com/fusesuphasin/go-fiber/app/utils/response"
	"github.com/fusesuphasin/go-fiber/app/utils/validation"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var pagePermission string = "permission"

// A UserController belong to the interface layer.
type PermissionController struct {
	Enforcer    *casbin.Enforcer
	Logger      interfaces.Logger
	Fiber       *fiber.App
	RoleService *service.RoleService
}

func NewPermissionController(logger interfaces.Logger, enforcer *casbin.Enforcer, app *fiber.App) *PermissionController {
	return &PermissionController{
		Enforcer: enforcer,
		Logger:   logger,
		Fiber:    app,
		RoleService: &service.RoleService{
			RoleRepository: repository.RoleRepository{
				Ctx: context.Background(),
			},
		},
	}
}

func (controller PermissionController) PermissionRouter() {
	// controller.Fiber.Group(pagePermission)
	controller.Fiber.Get("/permission/:id", middleware.CheckPermission(controller.Enforcer, pagePermission), controller.getListPermission)
	controller.Fiber.Post("/permission/:id" ,middleware.CheckPermission(controller.Enforcer, pagePermission), controller.updatePermission)
}

// List Permission
// @Summary List Permission
// @Description List Permission by role
// @Tags Permission
// @Param Authorization header string true "With the bearer started"
// @Accept application/json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /permission/:id [get]
func (controller PermissionController) getListPermission(c *fiber.Ctx) error {
	controller.Logger.LogAccess("%s %s %s\n", c.IP(), c.Method(), c.OriginalURL())
	id := c.Params("id")
	roleData := controller.RoleService.GetById(id)

	if roleData == nil {
		c.Status(404)
		return c.JSON(&response.ErrorResponse{
			Success: false,
			Message: "Data not found",
		})
	}
	
	getPolicy := controller.Enforcer.GetFilteredPolicy(0, roleData.Name)
	
	/* if err != nil {
		c.Status(422)
		return c.JSON(&response.ErrorResponse{
			Success: false,
			Message: "Silahkan periksa kembali",
		})
	}
 */
	return c.JSON(&response.SuccessResponse{
		Success: true,
		Message: "Successfully retrieved data",
		Data:    getPolicy,
	})

}

// Update Permission
// @Summary Update Permission
// @Description Update Permission by role
// @Tags Permission
// @Accept application/json
// @Param Authorization header string true "With the bearer started"
// @Param permission body rules.PermissionUpdate true "Permission"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /permission/:id [post]
func (controller PermissionController) updatePermission(c *fiber.Ctx) error {
	controller.Logger.LogAccess("%s %s %s\n", c.IP(), c.Method(), c.OriginalURL())

	var permission rules.PermissionUpdate
	err := c.BodyParser(&permission)

	if err != nil {
		c.Status(422)
		return c.JSON(&response.ErrorResponse{
			Success: false,
			Message: "Please check again",
		})
	}

	id := c.Params("id")

	roleData := controller.RoleService.GetById(id)

	if roleData == nil {
		c.Status(404)
		return c.JSON(&response.ErrorResponse{
			Success: false,
			Message: "Data not found",
		})
	}

	initval = validator.New()
	roleValidation(initval, *controller.RoleService)
	errVal := initval.Struct(permission)

	if errVal != nil {
		message := make(map[string]string)
		message["permission"] = "Make sure the data sent again"
		errorResponse := validation.ValidateRequest(errVal, message)
		return c.JSON(errorResponse)
	}

	controller.Enforcer.RemoveFilteredPolicy(0, roleData.Name)

	for _, s := range permission.Permission {
		controller.Enforcer.AddPolicy(roleData.Name, s.Page, s.Resource)
	}

	getPolicy := controller.Enforcer.GetFilteredPolicy(0, roleData.Name)

	if err != nil {
		c.Status(422)
		return c.JSON(&response.ErrorResponse{
			Success: false,
			Message: "Please check again",
		})
	}

	return c.JSON(&response.SuccessResponse{
		Success: true,
		Message: "Successfully retrieved data",
		Data:    getPolicy,
	})

}