package controller

import (
	"context"
	"time"

	"github.com/fusesuphasin/go-fiber/app/domain"
	"github.com/fusesuphasin/go-fiber/app/interfaces"
	"github.com/fusesuphasin/go-fiber/app/repository"
	"github.com/fusesuphasin/go-fiber/app/rules"

	"github.com/fusesuphasin/go-fiber/app/service"
	"github.com/fusesuphasin/go-fiber/app/utils/encrypt"
	"github.com/fusesuphasin/go-fiber/app/utils/jwt"
	"github.com/fusesuphasin/go-fiber/app/utils/response"
	"github.com/fusesuphasin/go-fiber/app/utils/validation"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var initval *validator.Validate

type Cookie struct {
	Authorization     string    `json:"Authorization"`
}

// A UserController belong to the interface layer.
type UserController struct {
	Userservice service.UserService
	Logger      interfaces.Logger
}

func NewUserController(logger interfaces.Logger) *UserController {
	return &UserController{
		Userservice: service.UserService{
			UserRepository: repository.UserRepository{
				// mongoHandler:   mongoHandler,
				Ctx: context.Background(),
			},
		},
		Logger: logger,
	}
}

// Register
// @Summary Register user
// @Description Register user
// @Tags Authentication
// @Accept application/json
// @Param register body rules.RegisterValidation true "Register"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /register [post]
func (controller UserController) Register(c *fiber.Ctx) error {
	//controller.Logger.LogAccess("%s %s %s\n", c.IP(), c.Method(), c.OriginalURL())
	var register *domain.User

	errRequest := c.BodyParser(&register)

	if errRequest != nil {
		c.Status(422)
		return c.JSON(&response.ErrorResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	initval = validator.New()
	
	
/* 	// register a custom function. The first parameter is tag custom in struct, and the second parameter is a custom function
	initval.RegisterValidation ("customervalidation", // register a custom function. The first parameter is tag custom in struct, and the second parameter is a custom function
	initval.RegisterValidation ("customervalidation", CustomerValidationFunc) 
)  */

	registerValidation(initval, controller.Userservice)
	
	// check validate
	errVal := initval.Struct(register)

	if errVal != nil {
		message := make(map[string]string)
		message["Username"] = "Username is duplicate"
		message["Password"] = "Password must more than 8 charactor"
		errorResponse := validation.ValidateRequest(errVal, message)
		c.Status(422)
		return c.JSON(&response.ErrorResponse{
			Success: false,
			Message: errorResponse,
		})
	}

	register.Password, _ = encrypt.CreateHash(register.Password, encrypt.DefaultParams)
	register.RoleID = "617ab3f465d6fda17eaee1af"

	data := controller.Userservice.CreateUser(register)

	token, errToken := jwt.CreateToken(data.Username)

	if errToken != nil {
		controller.Logger.LogError("%s", errToken)
	}

	persistToken := jwt.CreateAuth(controller.Userservice, data.Username, token)

	if persistToken != nil {
		controller.Logger.LogError("%s", errToken)
	}

	return c.JSON(&response.SuccessResponse{
		Success: true,
		Message: "Register Success"/* Berhasil mendaftar */,
		Data: &response.RegisterResponse{
			Name:     data.Name,
			Username: data.Username,
			Token:    token.AccessToken,
		},
	})

}

// Login
// @Summary Login user
// @Description Login user
// @Tags Authentication
// @Accept application/json
// @Param login body rules.LoginValidation true "Login"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /login [post]
func (controller UserController) Login(c *fiber.Ctx) error {

	controller.Logger.LogAccess("%s %s %s\n", c.IP(), c.Method(), c.OriginalURL())
	var login *rules.LoginValidation
	err := c.BodyParser(&login)
	if err != nil {
		_ = c.JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
	}

	initval = validator.New()
	loginValidation(initval, controller.Userservice)
	errVal := initval.Struct(login)

	if errVal != nil {
		message := make(map[string]string)
		message["Password"] = "Password must more than 6 charactor"
		errorResponse := validation.ValidateRequest(errVal, message)
		return c.JSON(errorResponse)
	}

	res := controller.Userservice.CheckUsername(login.Username)

	if res.Username == "" {
		c.Status(422)
		err = c.JSON(&fiber.Map{
			"success": false,
			"error":   "Data tidak ditemukan",
		})
		return err
	}

	check, _ := encrypt.ComparePasswordAndHash(login.Password, string(res.Password))

	if check {
		td, errToken := jwt.CreateToken(res.Username)
		if errToken != nil {
			controller.Logger.LogError("%s", errToken)
		}
		
		/* url := "http://127.0.0.1:3000/login"

		  // Create a Bearer string by appending string access token
		var bearer = "Bearer " + td.AccessToken
		
		// Create a new request using http
		req, err := http.NewRequest("GET", url, nil)
		fmt.Println("-------Req---------")
		fmt.Println(req)
		// add authorization header to the req
		req.Header.Add("Authorization", bearer)
		c.Set("FFFF", "text/plain")
		_ = err

	

		// Send req using http Client
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error on response.\n[ERROR] -", err)
		}
		defer resp.Body.Close()
		
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error while reading the response bytes:", err)
		}
		log.Println(string([]byte(body))) */
		cookie := fiber.Cookie{
			Name: "jwt",
			Value: td.AccessToken,
			Expires: time.Now().Add(time.Minute * 1),
			HTTPOnly: true,
		}

		c.Cookie(&cookie)

 		jwt.CreateAuth(controller.Userservice, res.Username, td)		

		return c.JSON(&response.SuccessResponse{
			Message: "login Success",
			Success: true,
			Data: &response.LoginResponse{
				Name:     res.Name,
				Username: res.Username,
				Token:    td.AccessToken,
			},
		})
	} else {
		c.Status(401)
		return c.JSON(&response.ErrorResponse{
			Success: false,
			Message: "Username/Password salah",
		})
	}

}

//If the key already exists, the previous validation function will be replaced. - this method is not thread-safe it is intended that these all be registered prior to any validation
func registerValidation(initval *validator.Validate, service service.UserService) {	
	
	initval.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		return IsValidUsername(service, fl.Field().String())
	})

	initval.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return IsValidPassword(service, fl.Field().String())
	})

}

func IsValidPassword(service service.UserService, input string) bool {
	return len(input) > 8
}

func IsValidUsername(service service.UserService, input string) bool {
	count := service.CheckUsernameCount(input)
	return count == int64(0)
}

func loginValidation(initval *validator.Validate, service service.UserService) {
	initval.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return IsValidPassword(service, fl.Field().String())
	})
}
