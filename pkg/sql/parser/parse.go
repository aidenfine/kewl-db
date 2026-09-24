package parser

import (
	"errors"
	"fmt"
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

	// this is hard coded and prob not correct we should fix this.
	tokensLen := len(tokens)
	if tokensLen < 2 {
		return ErrSyntaxError
	}

	// strip parens
	validParens, strippedTokens, err := stripAndCheckParens(tokens)
	if !validParens {
		return err
	}

	fmt.Println(strippedTokens)
	return nil
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

func stripAndCheckParens(tokens []lex.Token) (bool, []lex.Token, error) {
	validParens := []string{"(", ")", "[", "]"}

	parensMap := map[string]string{
		")": "(",
		"]": "[",
	}

	stack := []string{}
	strippedParensTokens := []lex.Token{}

	for _, v := range tokens {
		if slices.Contains(validParens, v.Value) {
			if opening, ok := parensMap[v.Value]; ok {
				if len(stack) == 0 {
					return false, tokens, ErrSyntaxError
				}

				if opening != stack[len(stack)-1] {
					return false, tokens, ErrSyntaxError
				}

				stack = stack[:len(stack)-1]
			} else {
				stack = append(stack, v.Value)
			}
		} else {
			strippedParensTokens = append(strippedParensTokens, v)
		}
	}

	if len(stack) != 0 {
		return false, tokens, ErrSyntaxError
	}

	return true, strippedParensTokens, nil
}
