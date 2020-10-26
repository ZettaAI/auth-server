package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomString(t *testing.T) {
	val, err := generateRandomString(12)
	if assert.NoError(t, err) {
		assert.Len(t, val, 16)
	}
}
