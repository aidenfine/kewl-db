package parser

import (
	"errors"
	"unicode"
	"unicode/utf8"

	"github.com/aidenfine/kewl-db/pkg/sql/lex"
)

var (
	ErrInvalidUTF8 = errors.New("UTF-8 characters are only allowed.")
	ErrInvalidChar = errors.New("Invalid character")
)

// TODO: finish this
func Parse(tokens []lex.Token) error {
	// initial token checks
	// We will run the optimizations here on the query (most likely)

	// check if the identifiers are invalid
	err := checkCharacters(tokens)
	return err

}

func checkCharacters(tokens []lex.Token) error {
	for _, v := range tokens {
		if !utf8.ValidString(v.Value) {
			return ErrInvalidUTF8
		}
		for _, c := range v.Value {
			if !(unicode.IsLetter(c) || unicode.IsNumber(c)) && c != '-' && c != '_' {
				return ErrInvalidChar
			}
		}
	}
	return nil
}
