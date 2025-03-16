package utils

import (
	"github.com/google/uuid"
)

func GenerateRandomFileName(key string) string {
	uniqueFileName := uuid.New().String() + "_" + key
	return uniqueFileName
}
