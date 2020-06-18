package store

import (
	"time"
	"unicode/utf8"
)

// Project represents an amadeus project.
// Only public fields are marshalled to JSON by default.
type Project struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Error     bool      `json:"-"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"createdAt"`
}

const (
	projNameMinLen = 3
	projNameMaxLen = 20
)

// validProjectNameRune checks if given project name rune is valid.
func validProjectNameRune(r rune) bool {
	if 'a' <= r && r <= 'z' {
		return true
	}
	if 'A' <= r && r <= 'Z' {
		return true
	}
	if '0' <= r && r <= '9' {
		return true
	}
	if r == '_' || r == '-' {
		return true
	}
	return false
}

// ValidProjectName checks if given project name is valid.
func ValidProjectName(projectName string) bool {
	length := utf8.RuneCountInString(projectName)
	if !(projNameMinLen <= length && length <= projNameMaxLen) {
		return false
	}

	for _, r := range projectName {
		if !validProjectNameRune(r) {
			return false
		}
	}

	return true
}
