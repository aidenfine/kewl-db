package parser

import (
	"errors"
	"slices"
	"unicode"
	"unicode/utf8"

	"github.com/aidenfine/kewl-db/pkg/sql/lex"
)

var (
	ErrInvalidUTF8 = errors.New("UTF-8 characters are only allowed.")
	ErrInvalidChar = errors.New("Invalid character")
	ErrSyntaxError = errors.New("Syntax error")
)

// TODO: finish this
// SQL statements must always start with one of the following
// SELECT, INSERT, UPDATE, DELETE, CREATE, ALTER, DROP, TRUNCATE, GRANT, REVOKE, BEGIN, COMMIT, ROLLBACK
func Parse(tokens []lex.Token) error {
	// initial token checks
	// We will run the optimizations here on the query (most likely)

	// based on query type call the validator for query type.

	// check if the identifiers are invalid
	err := checkCharacters(tokens)
	return err

}

func preCheck(tokens []lex.Token) error {
	// min len checks for tokens,
	// check if tokens end or start in a value that is not valid
	// validate parens, could use leetcode problem valid parentheses stack method for this.
	// After stripping all parens is the first token valid?

	tokensLen := len(tokens)
	if tokensLen < 2 {
		return ErrSyntaxError
	}

	firstToken := tokens[0]
	if firstToken.Value == "(" && tokensLen > 1 {
		firstToken = tokens[1]
	}
	if !slices.Contains(lex.ActiveQueryTypes, firstToken.Type) {
		return ErrSyntaxError

	}

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
