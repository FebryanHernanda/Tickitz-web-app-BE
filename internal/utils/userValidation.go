package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"

	"github.com/FebryanHernanda/Tickitz-web-app-BE/internal/models"
)

func IsValidEmail(user models.RegisterUser) error {
	reEmail, err := regexp.Compile(`^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,4}$`)
	if err != nil {
		return err
	}

	if isEmailMatched := reEmail.Match([]byte(user.Email)); !isEmailMatched {
		return errors.New("email format is incorrect, use format like: user@email.com")
	}

	return nil
}

func IsValidPassword(user models.RegisterUser) error {
	if len(user.Password) < 8 {
		return errors.New("password minimum 8 characters")
	}

	reUppercase := regexp.MustCompile(`[A-Z]`)
	if !reUppercase.MatchString(user.Password) {
		return errors.New("password harus mengandung huruf besar")
	}

	reLowercase := regexp.MustCompile(`[a-z]`)
	if !reLowercase.MatchString(user.Password) {
		return errors.New("password must contain capital letters")
	}

	reNumber := regexp.MustCompile(`[0-9]`)
	if !reNumber.MatchString(user.Password) {
		return errors.New("password must contain numbers")
	}

	reSpecial := regexp.MustCompile(`[#!?@$%^&*-]`)
	if !reSpecial.MatchString(user.Password) {
		return errors.New("password must contain special characters")
	}

	return nil
}

func GenerateVirtualAccount() (string, error) {
	result := "VA"
	for i := 0; i < 12; i++ {
		n := rand.Intn(10)
		result += fmt.Sprintf("%d", n)
	}
	return result, nil
}
