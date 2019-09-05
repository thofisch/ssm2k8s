package config

import (
	"strings"
)

var reservedWordList = []string{
	"managed",
}

var ReservedWords = strings.Join(reservedWordList, " ,")

func IsReservedWord(applicationName string) bool {
	for _, reservedWord := range reservedWordList {
		if reservedWord == applicationName {
			return true
		}
	}
	return false
}
