package utils

import (
	"github.com/mediocregopher/radix/v3"
	"strings"
)

const availableSymbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NumberToLexString(number int) string {
	var sb strings.Builder
	for number > 0 {
		newCharIndex := number % len(availableSymbols)
		sb.WriteByte(availableSymbols[newCharIndex])
		number -= newCharIndex
	}
	return sb.String()
}

func GetNewUrlInteger(pool *radix.Pool) int {
	var counter int
	_ = pool.Do(radix.Cmd(&counter, "INCR", "$counter"))
	return counter
}

func GetNewUrlString(pool *radix.Pool) string {
	return NumberToLexString(GetNewUrlInteger(pool))
}
