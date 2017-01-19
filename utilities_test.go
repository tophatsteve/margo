package margo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpperCase1stLetter(t *testing.T) {
	assert.Equal(t, isUppercase("Testing", 0), true, "First letter of 'Testing' should be uppercase")
}

func TestLowerCase2dnLetter(t *testing.T) {
	assert.Equal(t, isUppercase("Testing", 1), false, "Second letter of 'Testing' should be lowercase")
}
