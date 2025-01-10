package helper

import (
	"math/rand"
	"time"
)

type UtilInterfaceImpl struct{}

type UtilInterface interface {
	GenerateBlankToken() string
}

func NewUtil() UtilInterface {
	return &UtilInterfaceImpl{}
}

func (u *UtilInterfaceImpl) GenerateBlankToken() string {
	const charset = "abcdefgh0123456789"

	// Create a new random source and generator
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	result := make([]byte, 24)
	for i := range result {
		result[i] = charset[rng.Intn(len(charset))]
	}
	return string(result)
}
