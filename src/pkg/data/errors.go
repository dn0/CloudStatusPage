package data

import (
	"strings"
)

type InvalidInputError struct {
	Field string
}

type InvalidInputErrors []*InvalidInputError

func (e *InvalidInputError) Error() string {
	return "invalid input: " + e.Field
}

func (e InvalidInputErrors) Error() string {
	sb := new(strings.Builder)
	for _, err := range e {
		sb.WriteString(err.Error())
		sb.WriteString("\n")
	}
	return strings.TrimSpace(sb.String())
}
