package utils

import (
	"github.com/google/uuid"
)

func GenerateSKU() string {
	return "SKU-" + GenerateUUID()
}

func GenerateUUID() string {
	return uuid.New().String()
}
