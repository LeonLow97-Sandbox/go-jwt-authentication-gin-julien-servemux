package middleware

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func CompareHash(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func ReturnJson(w http.ResponseWriter, Code int, Message string) {
	jsonStatus := struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}{
		Message: Message,
		Code:    Code,
	}

	bs, err := json.Marshal(jsonStatus)
	if err != nil {
		log.Fatalln("Error in ReturnJson middleware", err)
	}
	io.WriteString(w, string(bs))
}

func InternalServerError(message string, err error) {
	log.Fatalln(message, err)
}