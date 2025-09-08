package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateQRCode() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	number := rand.Intn(90000) + 10000
	return "QR" + strconv.Itoa(number)
}
