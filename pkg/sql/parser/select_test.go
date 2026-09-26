package parser

import (
	"errors"
	"reflect"
	"testing"

	"github.com/aidenfine/kewl-db/pkg/sql/ast"
	"github.com/aidenfine/kewl-db/pkg/sql/lex"
)

func TestParseSelect(t *testing.T) {
	for _, tt := range []struct {
		query string
		want  *ast.SelectStatement
	}{
		{"SELECT * FROM users", &ast.SelectStatement{
			Columns: []ast.Expr{&ast.Star{}}, From: "users",
		}},
		{"select id, name from users where age>18", &ast.SelectStatement{
			Columns: []ast.Expr{&ast.ColumnRef{Name: "id"}, &ast.ColumnRef{Name: "name"}}, From: "users",
			Where: &ast.BinaryExpr{Left: &ast.ColumnRef{Name: "age"}, Op: ">", Right: &ast.IntLiteral{Value: 18}},
		}},
		{"SELECT id FROM users WHERE age = minimum_age", &ast.SelectStatement{
			Columns: []ast.Expr{&ast.ColumnRef{Name: "id"}}, From: "users",
			Where: &ast.BinaryExpr{Left: &ast.ColumnRef{Name: "age"}, Op: "=", Right: &ast.ColumnRef{Name: "minimum_age"}},
		}},
		{"SELECT id FROM users WHERE -1 < age", &ast.SelectStatement{
			Columns: []ast.Expr{&ast.ColumnRef{Name: "id"}}, From: "users",
			Where: &ast.BinaryExpr{Left: &ast.IntLiteral{Value: -1}, Op: "<", Right: &ast.ColumnRef{Name: "age"}},
		}},
	} {
		t.Run(tt.query, func(t *testing.T) {
			got, err := Parse(lex.Lex(tt.query))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestParseSelectSyntaxErrors(t *testing.T) {
	for _, query := range []string{
		"", "SELECT", "SELECT FROM users", "SELECT id", "SELECT id FROM",
		"SELECT id, FROM users", "SELECT id name FROM users", "SELECT *, id FROM users",
		"SELECT id FROM 42", "SELECT id FROM users WHERE", "SELECT id FROM users WHERE age",
		"SELECT id FROM users WHERE age >", "SELECT id FROM users WHERE age 18",
		"SELECT id FROM users WHERE age >= 18", "SELECT id FROM users extra",
		"SELECT id FROM users WHERE age > 18 AND age < 30",
		"SELECT id FROM users SELECT id FROM users",
	} {
		t.Run(query, func(t *testing.T) {
			got, err := Parse(lex.Lex(query))
			if !errors.Is(err, ErrSyntaxError) || got != nil {
				t.Fatalf("got (%v, %v), want (nil, syntax error)", got, err)
			}
		})
	}
}

func TestParseInvalidTokens(t *testing.T) {
	for _, query := range []string{"SELECT !!! FROM users", "SELECT id FROM user.name", "UPDATE users SET age = age + 1"} {
		if _, err := Parse(lex.Lex(query)); !errors.Is(err, ErrInvalidChar) {
			t.Errorf("%q: expected invalid character, got %v", query, err)
		}
	}
	tokens := lex.Lex("SELECT id FROM users WHERE age > 1")
	tokens[len(tokens)-1].Value = "9223372036854775808"
	if _, err := Parse(tokens); !errors.Is(err, ErrSyntaxError) {
		t.Errorf("expected integer overflow error, got %v", err)
	}
}
