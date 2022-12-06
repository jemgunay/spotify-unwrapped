package stats

import (
	"strings"
	"unicode"
)

// CountWordsInSentence counts the words in the given sentence and excludes boring words.
func CountWordsInSentence(sentence string, mapping Mapping) {
	for _, word := range strings.Split(sentence, " ") {
		trimmed := strings.TrimSpace(word)
		if len(trimmed) == 0 {
			continue
		}

		runes := []rune(trimmed)
		// validate that first character is valid
		firstChar := runes[0]
		if isCharInvalid(firstChar) {
			continue
		}

		// truncate weird words with trailing punctuation like closing brackets, etc
		if len(runes) > 1 {
			lastChar := runes[len(runes)-1]
			if isCharInvalid(lastChar) {
				runes = runes[:len(runes)-1]
			}
		}

		runes[0] = unicode.ToUpper(firstChar)
		word := string(runes)
		if shouldExclude(word) {
			continue
		}
		mapping.Push(word)
	}
}

func isCharInvalid(char rune) bool {
	return unicode.IsSymbol(char) || unicode.IsPunct(char) || unicode.IsNumber(char)
}

// filter boring words - words must be sentence-cased
var exclusionList = map[string]struct{}{
	"A":    {},
	"An":   {},
	"And":  {},
	"As":   {},
	"At":   {},
	"Be":   {},
	"For":  {},
	"In":   {},
	"Is":   {},
	"It":   {},
	"Of":   {},
	"On":   {},
	"So":   {},
	"The":  {},
	"To":   {},
	"With": {},
}

func shouldExclude(word string) bool {
	_, ok := exclusionList[word]
	return ok
}
