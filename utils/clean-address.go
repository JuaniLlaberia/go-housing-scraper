package utils

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func CleanAddress(address string, neighbour string) string {
	// Create new trimmed address
	newAddress := strings.ToLower(strings.TrimSpace(address))
	// Remove "Argentina" from string
	newAddress = strings.Replace(newAddress, "argentina", "", 1)
	// In case is "CABA"/"Ciudad Autonoma" we remove
	newAddress = regexp.MustCompile(`(?i)\b(capital federal|caba)\b`).ReplaceAllLiteralString(newAddress, "")
	// Remove extra symbols (".",",")
	newAddress = strings.ReplaceAll(newAddress, ".", "")
	newAddress = strings.ReplaceAll(newAddress, ",", "")
	// Remove neighbour
	newAddress = strings.Replace(newAddress, strings.ToLower(neighbour), "", 1)
	// Remove parenthesis and any data inside them
	newAddress = regexp.MustCompile(`\(.*\)`).ReplaceAllString(newAddress, "")

	// Capitalized
	caser := cases.Title(language.Spanish)
	newAddress = caser.String(newAddress)

	return newAddress
}
