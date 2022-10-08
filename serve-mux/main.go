package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/LeonLow97-Sandbox/go-jwt-login-authentication-methods/serve-mux/middleware"
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
	http.HandleFunc("/login", Login)
	http.HandleFunc("/logout", Logout)

	port := fmt.Sprintf(":%v", 8080)
	fmt.Println("Server is running on Port:", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("Error in running Port", err)
	}
}

func Login(w http.ResponseWriter, req *http.Request) {
	// database credentials
	dbUsername := "Leon"
	dbPassword := "$2a$10$YIofk.GOJ0jBlpjw1YzSKOQxcbN3bpyMRx82F4EDj/sBRpRkQlgDu"	// salt = 10, Password0!

	w.Header().Set("Content-Type", "application/json")
	var user UserCredentials

	bs, _ := io.ReadAll(req.Body)
	if err := json.Unmarshal(bs, &user); err != nil {
		middleware.InternalServerError("Error Occurred in Binding JSON", err)
		return
	}

	validatePassword := middleware.CompareHash(dbPassword, user.Password)
	if !validatePassword || dbUsername != user.Username {
		middleware.ReturnJson(w, http.StatusBadRequest, "Invalid Username or Password. Please try again")
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
		middleware.InternalServerError("Error in Signing Token", err)
		return
	}

	c := &http.Cookie {
		Name: "jwt-servemux",
		Value: token,
		MaxAge: 3600,
		Path: "/",
		Domain: "localhost",
		Secure: false,
		HttpOnly: true,
	}
	// Creating a Cookie with the Token (maxAge = 3600 = 1 hour)
	http.SetCookie(w, c)
	
	middleware.ReturnJson(w, 200, "Successfully Logged In!")
}

func Logout(w http.ResponseWriter, req *http.Request) {
	RetrieveIssuer(w, req)
	fmt.Println("Retrieved Issuer using w.Header().Get() in Logout Route:" , w.Header().Get("username"))

	w.Header().Set("Content-Type", "application/json")
	cookie, err := req.Cookie("jwt-servemux")
	if err != nil {
		middleware.InternalServerError("Error retrieving cookie.", err)
	}
	cookie.Value = ""

	cookie = &http.Cookie {
		Name: "jwt-servemux",
		Value: "",
		MaxAge: -1,
		Path: "",
		Domain: "",
		Secure: false,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	middleware.ReturnJson(w, 200, "Successfully Logged Out!")
}

func RetrieveIssuer(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("jwt-servemux")
	if err != nil {
		middleware.InternalServerError("Error retrieving cookie.", err)
	}

	// Parses the cookie jwt claims
	token, err := jwt.ParseWithClaims(cookie.Value, &jwt.RegisteredClaims{}, func(*jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		middleware.InternalServerError("Error parsing cookie", err)
		return
	}

	// To access the issuer, we need it to be type RegisteredClaims
	claims := token.Claims.(*jwt.RegisteredClaims)

	fmt.Println("Retrieved Issuer using claims.Issuer:", claims.Issuer)
	w.Header().Set("username", claims.Issuer)
	fmt.Println("Retrieved Issuer using w.Header().Get():" , w.Header().Get("username"))
}