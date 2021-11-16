package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/fusesuphasin/go-fiber/app/service"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)



type AccessDetails struct {
	AccessUuid string
	UserName     string
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string `json:"-"`
	RefreshUuid  string `json:"-"`
	AtExpires    int64  `json:"-"`
	RtExpires    int64  `json:"-"`
}

func CreateToken(userName string) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewString()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = td.AccessUuid + "++" + userName
	
	var err error
	//Creating Access Token
	// os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_name"] = userName
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	// os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_name"] = userName
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func CreateAuth(service service.UserService, userName string, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := service.InsertToken(td.AccessUuid, userName, at.Sub(now))
	service.InsertToken(td.RefreshUuid, userName, rt.Sub(now))

	return errAccess
}

func ExtractTokenMetadata(r *fiber.Ctx) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	fmt.Println("___claims____")
	fmt.Println(claims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		//userName, err := strconv.ParseUint(fmt.Sprintf("%.s", claims["user_name"]), 10, 64)
		userName := fmt.Sprintf("%s", claims["user_name"])
		/* if err != nil {
			return nil, err
		} */
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserName:     userName,
		}, nil
	}
	return nil, err
}

func VerifyToken(r *fiber.Ctx) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("___tokenString____")
	fmt.Println(tokenString)
	fmt.Println("___Token____")
	fmt.Println(token)
	return token, nil
}

func ExtractToken(r *fiber.Ctx) string {
	cookie := r.Cookies("jwt")
	if cookie != "" {
		return cookie
	}
	/* bearToken := r.Get("Authorization")
	strArr := strings.Split(bearToken, " ")

	if len(strArr) == 2 {
		return strArr[1]
	} */
	return ""
}

func FetchAuth(service service.UserService, authD *AccessDetails) (string, string, error) {
	userName, _ := service.FetchToken(authD.AccessUuid)
	if authD.UserName != userName {
		return "", "", fmt.Errorf("userName mismatch")
	}
	
	return userName, authD.UserName, nil
}
