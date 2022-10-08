package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/LeonLow97-Sandbox/go-jwt-login-authentication-methods/gin-gonic/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	SecretKey = "qwerty123"
)

func main() {
	router := gin.Default()
	router.POST("/login", Login)
	router.GET("/logout", Logout)

	port := fmt.Sprintf(":%v", 8080)
	log.Fatalln(router.Run(port))
}

func Login(c *gin.Context) {
	// database credentials
	dbUsername := "Leon"
	dbPassword := "$2a$10$YIofk.GOJ0jBlpjw1YzSKOQxcbN3bpyMRx82F4EDj/sBRpRkQlgDu"	// salt = 10, Password0!

	var user UserCredentials

	if err := c.BindJSON(&user); err != nil {
		middleware.InternalServerError(c, "Error Occurred in Binding JSON", err)
		return
	}

	validatePassword := middleware.CompareHash(dbPassword, user.Password)
	if !validatePassword || dbUsername != user.Username {
		middleware.ErrorHandler(c, http.StatusBadRequest, "Invalid Username or Password. Please try again")
		return
	}

	// Generate Token with Claims
	tokenExpireTime := time.Now().Add(1 * time.Hour) // 1 hour expiry time from now
	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims {
		Issuer : dbUsername,
		ExpiresAt: jwt.NewNumericDate(tokenExpireTime),	
	})

	// Used to validate that the token is trustworthy and has not been tampered with a secret key (usually stored in .env)
	token, err := generateToken.SignedString([]byte(SecretKey))
	if err != nil {
		middleware.InternalServerError(c, "Error in Signing Token", err)
		return
	}

	// Creating a Cookie with the Token (maxAge = 3600 = 1 hour)
	c.SetCookie("jwt-gin", token, 3600, "/", "localhost", false, true)
	
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message": "Successfully Logged In!",
	})
}

func Logout (c *gin.Context) {
	RetrieveIssuer(c)
	fmt.Println("Retrieved Issuer using c.GetString() in Logout Route:", c.GetString("username"))

	// Eliminates the cookie with maxAge = -1
	c.SetCookie("jwt-gin", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"message": "Successfully Logged Out!",
	})
}

// Retrieve the Issuer from Cookie
func RetrieveIssuer(c *gin.Context) {
	// Returns existing cookie
	cookie, err := c.Cookie("jwt-gin")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"code": http.StatusBadRequest, "message": "Invalid Credentials"})
		return
	}
	
	// Parses the cookie jwt claims
	token, err := jwt.ParseWithClaims(cookie, &jwt.RegisteredClaims{}, func(*jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		middleware.InternalServerError(c, "Error parsing cookie", err)
		return
	}

	// To access the issuer, we need it to be type RegisteredClaims
	claims := token.Claims.(*jwt.RegisteredClaims)

	// Retrieve issuer from Claims
	// and Set issue as "username" to access in backend
	fmt.Println("Retrieved Issuer using claims.Issuer:", claims.Issuer)
	c.Set("username", claims.Issuer)
	fmt.Println("Retrieved Issuer using c.GetString():" , c.GetString("username"))
}