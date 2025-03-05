package util

import (
	"crypto/rand"
	"log"
	"strconv"
)

func GenerateBankAccountNumber() string {
	bytes := make([]byte, 13)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatalf("Failed to generate random number: %v", err)
	}

	randomNumber := ""
	for _, b := range bytes {
		randomNumber += strconv.Itoa(int(b) % 10)
	}

	accountNumber := "820" + "123" + randomNumber

	checkDigit := 0
	for _, char := range accountNumber {
		digit := int(char - '0')
		checkDigit += digit
	}

	checkDigit = checkDigit % 10
	accountNumber += strconv.Itoa(checkDigit)

	return accountNumber
}
