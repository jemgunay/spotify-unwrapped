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
		char := runes[0]
		if unicode.IsSymbol(char) || unicode.IsPunct(char) || unicode.IsNumber(char) {
			continue
		}

		// truncate weird trailing punctuation
		if len(trimmed) < 1 {
			lastChar := runes[len(runes)-1]
			if unicode.IsSymbol(lastChar) || unicode.IsPunct(lastChar) || unicode.IsNumber(lastChar) {
				runes = runes[:len(runes)-1]
			}
		}

		runes[0] = unicode.ToUpper(char)
		word := string(runes)
		if shouldExclude(word) {
			continue
		}
		mapping.Push(word)
	}
}

// words must be sentence-cased
var exclusionList = map[string]struct{}{
	"A":   {},
	"An":  {},
	"As":  {},
	"At":  {},
	"In":  {},
	"It":  {},
	"Of":  {},
	"On":  {},
	"The": {},
	"To":  {},
}

func shouldExclude(word string) bool {
	_, ok := exclusionList[word]
	return ok
}
