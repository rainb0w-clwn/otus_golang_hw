package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func isSlashSymbol(s rune) bool {
	return string(s) == "\\"
}

func writeSymbol(r *strings.Builder, s *rune, c int) {
	r.WriteString(strings.Repeat(string(*s), c))
	*s = 0
}

func Unpack(input string) (string, error) {
	var result strings.Builder
	var pSymbol rune
	var pEscaped bool
	for _, cSymbol := range input {
		counter, err := strconv.Atoi(string(cSymbol))
		switch {
		case err != nil:
			if pEscaped {
				if isSlashSymbol(cSymbol) {
					pSymbol = cSymbol
					pEscaped = false
					continue
				}
				return "", ErrInvalidString // "заэкранировать можно только цифру или слэш"
			}
			if pSymbol != 0 {
				writeSymbol(&result, &pSymbol, 1)
			}
			pEscaped = isSlashSymbol(cSymbol)
			pSymbol = cSymbol
		case pEscaped:
			pSymbol = cSymbol
			pEscaped = false
		case pSymbol == 0:
			return "", ErrInvalidString
		default:
			writeSymbol(&result, &pSymbol, counter)
		}
	}
	if pSymbol != 0 && !pEscaped {
		writeSymbol(&result, &pSymbol, 1)
	}
	return result.String(), nil
}
