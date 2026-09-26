package parser

import (
	"errors"
	"reflect"
	"testing"

	"github.com/aidenfine/kewl-db/pkg/sql/ast"
	"github.com/aidenfine/kewl-db/pkg/sql/lex"
)

func TestParseStatements(t *testing.T) {
	for _, tt := range []struct {
		query string
		want  ast.Statement
	}{
		{"INSERT INTO users (id, name) VALUES (1, 'Alice'), (2, 'O''Brien');", &ast.InsertStatement{
			Table: "users", Columns: []string{"id", "name"},
			Values: [][]ast.Expr{
				{&ast.IntLiteral{Value: 1}, &ast.StringLiteral{Value: "Alice"}},
				{&ast.IntLiteral{Value: 2}, &ast.StringLiteral{Value: "O'Brien"}},
			},
		}},
		{"insert into users values (-2, ''), (3, 'hi, (there); > =');", &ast.InsertStatement{
			Table: "users", Values: [][]ast.Expr{
				{&ast.IntLiteral{Value: -2}, &ast.StringLiteral{Value: ""}},
				{&ast.IntLiteral{Value: 3}, &ast.StringLiteral{Value: "hi, (there); > ="}},
			},
		}},
		{"UPDATE users SET name = 'Alice', age = other_age WHERE id = 1;", &ast.UpdateStatement{
			Table: "users", Assignments: []ast.Assignment{
				{Column: "name", Value: &ast.StringLiteral{Value: "Alice"}},
				{Column: "age", Value: &ast.ColumnRef{Name: "other_age"}},
			}, Where: &ast.BinaryExpr{Left: &ast.ColumnRef{Name: "id"}, Op: "=", Right: &ast.IntLiteral{Value: 1}},
		}},
		{"UPDATE users SET age = 18", &ast.UpdateStatement{
			Table: "users", Assignments: []ast.Assignment{{Column: "age", Value: &ast.IntLiteral{Value: 18}}},
		}},
		{"DELETE FROM users WHERE name = 'Alice';", &ast.DeleteStatement{
			Table: "users", Where: &ast.BinaryExpr{Left: &ast.ColumnRef{Name: "name"}, Op: "=", Right: &ast.StringLiteral{Value: "Alice"}},
		}},
		{"DELETE FROM users", &ast.DeleteStatement{Table: "users"}},
		{"create database MyDB;", &ast.CreateDatabaseStatement{Name: "MyDB"}},
		{"DROP DATABASE MyDB", &ast.DropDatabaseStatement{Name: "MyDB"}},
		{"SELECT * FROM users WHERE name = 'Alice';", &ast.SelectStatement{
			Columns: []ast.Expr{&ast.Star{}}, From: "users",
			Where: &ast.BinaryExpr{Left: &ast.ColumnRef{Name: "name"}, Op: "=", Right: &ast.StringLiteral{Value: "Alice"}},
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

func TestStatementSyntaxErrors(t *testing.T) {
	for _, query := range []string{
		"INSERT", "INSERT users VALUES (1)", "INSERT INTO", "INSERT INTO users",
		"INSERT INTO users () VALUES (1)", "INSERT INTO users (id,) VALUES (1)",
		"INSERT INTO users (id VALUES (1)", "INSERT INTO users VALUES",
		"INSERT INTO users VALUES ()", "INSERT INTO users VALUES (1,)",
		"INSERT INTO users VALUES (1", "INSERT INTO users VALUES (1),",
		"INSERT INTO users VALUES (1) (2)", "INSERT INTO users VALUES (unquoted)",
		"INSERT INTO users (id, name) VALUES (1)", "INSERT INTO users VALUES (1), (2, 3)",
		"INSERT INTO users VALUES (9223372036854775808)",
		"INSERT INTO users VALUES ('unterminated)",
		"UPDATE", "UPDATE users", "UPDATE users SET", "UPDATE users name = 'Alice'",
		"UPDATE users SET name", "UPDATE users SET name =", "UPDATE users SET name > 1",
		"UPDATE users SET name = 1,", "UPDATE users SET name = 1 age = 2",
		"UPDATE users SET name = 1 WHERE",
		"DELETE", "DELETE users", "DELETE FROM", "DELETE FROM users WHERE id =",
		"CREATE", "CREATE DATABASE", "CREATE TABLE users", "CREATE DATABASE db extra",
		"DROP", "DROP DATABASE", "DROP TABLE users", "DROP DATABASE db extra",
		"DELETE FROM users;;", "DELETE FROM users; DROP DATABASE db",
		"SELECT * FROM users WHERE name = 'a' 'b'",
	} {
		t.Run(query, func(t *testing.T) {
			got, err := Parse(lex.Lex(query))
			if !errors.Is(err, ErrSyntaxError) || got != nil {
				t.Fatalf("got (%v, %v), want (nil, syntax error)", got, err)
			}
		})
	}
}

func FuzzParse(f *testing.F) {
	for _, query := range []string{"", "INSERT INTO t VALUES ('a', 1)", "UPDATE t SET x = 1 WHERE y > 2", "DELETE FROM t", "CREATE DATABASE db", "SELECT * FROM t"} {
		f.Add(query)
	}
	f.Fuzz(func(t *testing.T, query string) {
		stmt, err := Parse(lex.Lex(query))
		if err != nil && stmt != nil {
			t.Fatal("failed parse returned an AST")
		}
		if err == nil && stmt == nil {
			t.Fatal("successful parse returned no AST")
		}
	})
}
