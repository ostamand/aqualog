package util

import (
	"math/rand"
	"strings"
)

func GenerateRandomString(size int) string {
	chars := "abcdefghijklmnopqrtuvwxyz"
	n := len(chars)
	var sb strings.Builder
	for i := 0; i < size; i++ {
		sb.WriteByte(chars[rand.Intn(n)])
	}
	return sb.String()
}
