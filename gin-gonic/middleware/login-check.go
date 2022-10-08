package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CompareHash(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func ErrorHandler(c *gin.Context, code int, message interface{}) {
	c.JSON(code, gin.H{"code": code, "message": message})
}

func InternalServerError(c *gin.Context, message string, err error) {
	log.Fatalln(message, err)
	c.JSON(http.StatusInternalServerError, message)
}