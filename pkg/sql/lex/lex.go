package lex

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// SQL examples
// SELECT id, name FROM users WHERE age > 18
// After lex
// [SELECT, ID(id), COMMA, ID(name), FROM, ID(users), WHERE, ID(age), GT(>), INT(18)]

type tokenType int

// Defines the token types for the lex. With a defined token struct we could simply
// use the TokenType like this.
//
//	type Token struct {
//		Type TokenType
//		Value string
//	}
//
// This makes it super easy to be used like
//
//	token := Token{
//		Type: SELECT,
//		Value: "SELECT",
//	}
//
// and for a query like `SELECT users FROM accounts WHERE id > 10`
//
// We can easily id which table's to query from because Token will be something like...
//
//	Token{
//		Type: ID
//		Value: "accounts",
//	}
const (
	SELECT tokenType = iota
	FROM
	WHERE

	ID
	INT
	STRING

	EQ
	GT
	LT

	COMMA
	STAR
)

type Token struct {
	Type  tokenType
	Value string
}

// the lexer should take in a sql string and return a list of `Tokens`

// lexer should go char by char until it hits a whitespace or something that indicates a token has ended.
// So for "SELECT" it would go by each char and since SELECT is not an identifier we can assume that is a finished token
//
// Now for identifiers if we take the sql query `SELECT users FROM accounts WHERE id > 10`
// Our identifiers (ID) would be users, accounts, id, 10.
//

// lex will take in a sql query and return a []Token. As of now this is very loosely created and will really take anything in.
// We should make this much more strict in the future. One example is if we take `SELECT !!!` the lex will
// Just assume !!! is valid.
// TOOD: make this more strict
func Lex(sql string) []Token {
	currStr := ""
	tokens := []Token{}

	flush := func() {
		if currStr == "" {
			return
		}

		tokens = append(tokens, Token{Type: getTokenType(currStr), Value: currStr})

		currStr = ""
	}
	for _, c := range sql {
		if unicode.IsSpace(c) {
			flush()
			continue
		}
		switch c {
		case '=':
			flush()
			tokens = append(tokens, Token{Type: EQ, Value: string(c)})
		case '>':
			flush()
			tokens = append(tokens, Token{Type: GT, Value: string(c)})
		case '<':
			flush()
			tokens = append(tokens, Token{Type: LT, Value: string(c)})
		case ',':
			flush()
			tokens = append(tokens, Token{Type: COMMA, Value: string(c)})
		case '*':
			flush()
			tokens = append(tokens, Token{Type: STAR, Value: string(c)})
		default:
			currStr += string(c)
		}

	}
	flush()

	return tokens
}
func CheckCharacters(tokens []Token) error {
	for _, v := range tokens {
		if !utf8.ValidString(v.Value) {
			return errors.New("UTF-8 characters are only allowed.")
		}
		for _, c := range v.Value {
			if !(unicode.IsLetter(c) || unicode.IsNumber(c)) && c != '-' && c != '_' {
				return errors.New("Invalid character")
			}
		}
	}
	return nil
}

// id will be default if none is found (this may be bad assumtion?)
func getTokenType(str string) tokenType {
	switch strings.ToUpper(str) {
	case "SELECT":
		return SELECT
	case "FROM":
		return FROM
	case "WHERE":
		return WHERE
	default:
		if _, err := strconv.Atoi(str); err == nil {
			return INT
		}
		return ID
	}

}
