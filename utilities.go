package margo

import (
	"math/rand"
	"strings"
	"time"
)

func randomNumber(max int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Intn(max)
}

func isUppercase(s string, pos int) bool {
	c := string(s[pos])
	if c != strings.ToUpper(c) {
		return false
	}

	return true
}
