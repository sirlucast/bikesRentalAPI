package helpers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func WriteJSON(rw http.ResponseWriter, status int, data interface{}) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	_, err = rw.Write(js)
	if err != nil {
		return err
	}
	return nil
}

// ParseBody parses the body of a request and returns a map of strings
func ParseBody(bodyReadCloser io.ReadCloser) ([]byte, error) {
	body, err := io.ReadAll(bodyReadCloser)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %v", err)
	}
	defer bodyReadCloser.Close()
	return body, nil
}

// getHashPassword calls bcrypts function to give us the hashed password
func GetHashPassword(password string) (string, error) {
	bytePassword := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return "", err
	}
	return string(hash), nil
}

// base64Decode decodes a base64 string
func Base64Decode(str string) (string, bool) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", true
	}
	return string(data), false
}

// CheckPassword compares the hashed password with a string password
func CheckPassword(hashPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}
