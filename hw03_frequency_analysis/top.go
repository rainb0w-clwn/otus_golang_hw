package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var regular = regexp.MustCompile(`(\p{L}.*\p{L}+)`)

func GetSafeWord(word string) string {
	if word == "-" {
		return ""
	}
	safeWord := regular.FindString(word)
	if len(safeWord) > 0 { // рег.выражение требует минимум 2 "буквы"
		word = safeWord
	}
	// один знак препинания (не тире) или последовательность только из знаков препинания удовлетворяет условию
	return strings.ToLower(word)
}

func Top10(input string) []string {
	var words []string
	countMap := make(map[string]uint)
	for _, word := range strings.Fields(input) {
		word = GetSafeWord(word)
		if len(word) > 0 {
			if _, exists := countMap[word]; !exists {
				words = append(words, word)
			}
			countMap[word]++
		}
	}
	sort.Slice(words, func(i, j int) bool {
		if countMap[words[i]] > countMap[words[j]] {
			return true
		} else if countMap[words[i]] == countMap[words[j]] {
			return words[i] < words[j]
		}
		return false
	})
	return words[:min(10, len(words))]
}
