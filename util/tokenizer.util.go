package util

import (
	"strings"
	"unicode"
)

func Tokenize(text string) []string {
	text = strings.ToLower(text)
	estimatedTokens := len(text) / 6
	if estimatedTokens < 8 {
		estimatedTokens = 8
	}
	tokens := make([]string, 0, estimatedTokens)
	var builder strings.Builder
	builder.Grow(32)

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			builder.WriteRune(r)
		} else if builder.Len() > 0 {
			if token := builder.String(); len(token) > 1 {
				tokens = append(tokens, token)
			}
			builder.Reset()
		}
	}

	if builder.Len() > 1 {
		tokens = append(tokens, builder.String())
	}

	return tokens
}

func GenerateTokens(fields []string, tags []string) []string {
	var builder strings.Builder
	for _, field := range fields {
		builder.WriteString(field)
		builder.WriteByte(' ')
	}
	builder.WriteString(strings.Join(tags, " "))
	return Tokenize(builder.String())
}
