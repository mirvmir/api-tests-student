package utils

import "unicode"

func MixCase(text string) string {
	result := []rune(text)
	upper := false

	for i, char := range result {
		if !unicode.IsLetter(char) {
			continue
		}

		if upper {
			result[i] = unicode.ToUpper(char)
		} else {
			result[i] = unicode.ToLower(char)
		}
		upper = !upper
	}

	return string(result)
}
