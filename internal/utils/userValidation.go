package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
)

func IsValidEmail(email string) error {
	reEmail, err := regexp.Compile(`^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,4}$`)
	if err != nil {
		return err
	}

	if isEmailMatched := reEmail.Match([]byte(email)); !isEmailMatched {
		return errors.New("email format is incorrect, use format like: user@email.com")
	}

	return nil
}

func IsValidPassword(password string) error {
	if len(password) < 8 {
		return errors.New("password minimum 8 characters")
	}

	reUppercase := regexp.MustCompile(`[A-Z]`)
	if !reUppercase.MatchString(password) {
		return errors.New("password harus mengandung huruf besar")
	}

	reLowercase := regexp.MustCompile(`[a-z]`)
	if !reLowercase.MatchString(password) {
		return errors.New("password must contain capital letters")
	}

	reNumber := regexp.MustCompile(`[0-9]`)
	if !reNumber.MatchString(password) {
		return errors.New("password must contain numbers")
	}

	reSpecial := regexp.MustCompile(`[#!?@$%^&*-]`)
	if !reSpecial.MatchString(password) {
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
