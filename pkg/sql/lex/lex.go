package lex

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// SQL examples
// SELECT id, name FROM users WHERE age > 18
// After lex
// [SELECT, ID(id), COMMA, ID(name), FROM, ID(users), WHERE, ID(age), GT(>), INT(18)]

type TokenType int

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
	SELECT TokenType = iota
	FROM
	WHERE

	INSERT
	INTO
	VALUES

	DELETE
	UPDATE
	DROP

	ID
	INT
	STRING

	EQ
	GT
	LT

	COMMA
	STAR
	CREATE
	DATABASE
	SET
	LPAREN
	RPAREN
	SEMICOLON
	INVALID
)

var ActiveQueryTypes = []TokenType{SELECT, INSERT, DELETE, UPDATE, DROP, CREATE}

type Token struct {
	Type  TokenType
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
	if !utf8.ValidString(sql) {
		return []Token{{Type: INVALID, Value: sql}}
	}
	currStr := ""
	tokens := []Token{}

	flush := func() {
		if currStr == "" {
			return
		}

		tokens = append(tokens, Token{Type: getTokenType(currStr), Value: currStr})

		currStr = ""
	}
	runes := []rune(sql)
	for i := 0; i < len(runes); i++ {
		c := runes[i]
		if unicode.IsSpace(c) {
			flush()
			continue
		}
		switch c {
		case '\'':
			flush()
			var value strings.Builder
			closed := false
			for i++; i < len(runes); i++ {
				if runes[i] == '\'' {
					if i+1 < len(runes) && runes[i+1] == '\'' {
						value.WriteRune('\'')
						i++
						continue
					}
					closed = true
					break
				}
				value.WriteRune(runes[i])
			}
			kind := STRING
			if !closed {
				kind = INVALID
			}
			tokens = append(tokens, Token{Type: kind, Value: value.String()})
		case '(':
			flush()
			tokens = append(tokens, Token{Type: LPAREN, Value: "("})
		case ')':
			flush()
			tokens = append(tokens, Token{Type: RPAREN, Value: ")"})
		case ';':
			flush()
			tokens = append(tokens, Token{Type: SEMICOLON, Value: ";"})
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

// id will be default if none is found (this may be bad assumtion?)
func getTokenType(str string) TokenType {
	switch strings.ToUpper(str) {
	case "SELECT":
		return SELECT
	case "FROM":
		return FROM
	case "WHERE":
		return WHERE
	case "INSERT":
		return INSERT
	case "VALUES":
		return VALUES
	case "INTO":
		return INTO
	case "DELETE":
		return DELETE
	case "DROP":
		return DROP
	case "UPDATE":
		return UPDATE
	case "CREATE":
		return CREATE
	case "DATABASE":
		return DATABASE
	case "SET":
		return SET
	default:
		if isInteger(str) {
			return INT
		}
		return ID
	}

}

// Classify integers by spelling so overflow is reported by the parser,
// rather than accidentally treating a large numeric literal as an identifier.
func isInteger(s string) bool {
	if len(s) > 0 && (s[0] == '-' || s[0] == '+') {
		s = s[1:]
	}
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
