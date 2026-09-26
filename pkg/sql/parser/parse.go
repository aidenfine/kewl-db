package parser

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/aidenfine/kewl-db/pkg/sql/ast"
	"github.com/aidenfine/kewl-db/pkg/sql/lex"
)

var (
	ErrInvalidUTF8 = errors.New("UTF-8 characters are only allowed.")
	ErrInvalidChar = errors.New("Invalid character")
	ErrSyntaxError = errors.New("Syntax error")
)

// Parse parses one supported statement, optionally followed by a semicolon.
// It validates syntax, but does not check names or types against a schema.
func Parse(tokens []lex.Token) (ast.Statement, error) {
	if err := checkCharacters(tokens); err != nil {
		return nil, err
	}
	p := &Parser{tokens: tokens}
	var stmt ast.Statement
	var err error
	if len(tokens) == 0 {
		return nil, p.syntaxError("statement")
	}
	switch tokens[0].Type {
	case lex.SELECT:
		stmt, err = p.parseSelect()
	case lex.INSERT:
		stmt, err = p.parseInsert()
	case lex.UPDATE:
		stmt, err = p.parseUpdate()
	case lex.DELETE:
		stmt, err = p.parseDelete()
	case lex.CREATE, lex.DROP:
		stmt, err = p.parseDatabase()
	default:
		return nil, p.syntaxError("SELECT, INSERT, UPDATE, DELETE, CREATE, or DROP")
	}
	if err != nil {
		return nil, err
	}
	p.match(lex.SEMICOLON)
	if p.pos != len(tokens) {
		return nil, p.syntaxError("end of query")
	}
	return stmt, nil
}

type Parser struct {
	tokens []lex.Token
	pos    int
}

func (p *Parser) match(kind lex.TokenType) bool {
	if p.pos < len(p.tokens) && p.tokens[p.pos].Type == kind {
		p.pos++
		return true
	}
	return false
}

func (p *Parser) expect(kind lex.TokenType, description string) (lex.Token, error) {
	if !p.match(kind) {
		return lex.Token{}, p.syntaxError(description)
	}
	return p.tokens[p.pos-1], nil
}

func (p *Parser) syntaxError(expected string) error {
	if p.pos == len(p.tokens) {
		return fmt.Errorf("%w: expected %s at end of query", ErrSyntaxError, expected)
	}
	return fmt.Errorf("%w: expected %s at token %d, got %q", ErrSyntaxError, expected, p.pos+1, p.tokens[p.pos].Value)
}

func (p *Parser) parseSelect() (*ast.SelectStatement, error) {
	if _, err := p.expect(lex.SELECT, "SELECT"); err != nil {
		return nil, err
	}
	stmt := &ast.SelectStatement{}
	if p.match(lex.STAR) {
		stmt.Columns = []ast.Expr{&ast.Star{}}
	} else {
		for {
			column, err := p.expect(lex.ID, "column name")
			if err != nil {
				return nil, err
			}
			stmt.Columns = append(stmt.Columns, &ast.ColumnRef{Name: column.Value})
			if !p.match(lex.COMMA) {
				break
			}
		}
	}
	if _, err := p.expect(lex.FROM, "FROM"); err != nil {
		return nil, err
	}
	table, err := p.expect(lex.ID, "table name")
	if err != nil {
		return nil, err
	}
	stmt.From = table.Value
	stmt.Where, err = p.parseOptionalWhere()
	if err != nil {
		return nil, err
	}
	return stmt, nil
}

func (p *Parser) parseComparison() (ast.Expr, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	var op string
	switch {
	case p.match(lex.EQ):
		op = "="
	case p.match(lex.GT):
		op = ">"
	case p.match(lex.LT):
		op = "<"
	default:
		return nil, p.syntaxError("=, >, or <")
	}
	right, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	return &ast.BinaryExpr{Left: left, Op: op, Right: right}, nil
}

func (p *Parser) parsePrimary() (ast.Expr, error) {
	if p.match(lex.ID) {
		return &ast.ColumnRef{Name: p.tokens[p.pos-1].Value}, nil
	}
	return p.parseLiteral()
}

func (p *Parser) parseLiteral() (ast.Expr, error) {
	if p.match(lex.STRING) {
		return &ast.StringLiteral{Value: p.tokens[p.pos-1].Value}, nil
	}
	if p.match(lex.INT) {
		token := p.tokens[p.pos-1]
		value, err := strconv.ParseInt(token.Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: invalid integer %q", ErrSyntaxError, token.Value)
		}
		return &ast.IntLiteral{Value: value}, nil
	}
	return nil, p.syntaxError("integer or quoted string")
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
		if v.Type == lex.INVALID {
			return fmt.Errorf("%w: invalid token or unterminated string", ErrSyntaxError)
		}
		if !utf8.ValidString(v.Value) {
			return ErrInvalidUTF8
		}
		if v.Type != lex.ID {
			continue
		}
		if v.Value == "" {
			return ErrInvalidChar
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
