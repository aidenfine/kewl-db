package lex_test

import (
	"reflect"
	"testing"

	"github.com/aidenfine/kewl-db/pkg/sql/lex"
)

func TestLex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []lex.Token
	}{
		{
			name:  "select",
			input: "SELECT",
			expected: []lex.Token{
				{Type: lex.SELECT, Value: "SELECT"},
			},
		},
		{
			name:  "select from where",
			input: "SELECT FROM WHERE",
			expected: []lex.Token{
				{Type: lex.SELECT, Value: "SELECT"},
				{Type: lex.FROM, Value: "FROM"},
				{Type: lex.WHERE, Value: "WHERE"},
			},
		},
		{
			name:  "identifiers",
			input: "users accounts id",
			expected: []lex.Token{
				{Type: lex.ID, Value: "users"},
				{Type: lex.ID, Value: "accounts"},
				{Type: lex.ID, Value: "id"},
			},
		},
		{
			name:  "integer",
			input: "10",
			expected: []lex.Token{
				{Type: lex.INT, Value: "10"},
			},
		},
		{
			name:  "comparison operators",
			input: "= > <",
			expected: []lex.Token{
				{Type: lex.EQ, Value: "="},
				{Type: lex.GT, Value: ">"},
				{Type: lex.LT, Value: "<"},
			},
		},
		{
			name:  "comma and star",
			input: ", *",
			expected: []lex.Token{
				{Type: lex.COMMA, Value: ","},
				{Type: lex.STAR, Value: "*"},
			},
		},
		{
			name:  "simple select query",
			input: "SELECT users FROM accounts",
			expected: []lex.Token{
				{Type: lex.SELECT, Value: "SELECT"},
				{Type: lex.ID, Value: "users"},
				{Type: lex.FROM, Value: "FROM"},
				{Type: lex.ID, Value: "accounts"},
			},
		},
		{
			name:  "select with where",
			input: "SELECT users FROM accounts WHERE id > 10",
			expected: []lex.Token{
				{Type: lex.SELECT, Value: "SELECT"},
				{Type: lex.ID, Value: "users"},
				{Type: lex.FROM, Value: "FROM"},
				{Type: lex.ID, Value: "accounts"},
				{Type: lex.WHERE, Value: "WHERE"},
				{Type: lex.ID, Value: "id"},
				{Type: lex.GT, Value: ">"},
				{Type: lex.INT, Value: "10"},
			},
		},
		{
			name:  "operators without whitespace",
			input: "id>10",
			expected: []lex.Token{
				{Type: lex.ID, Value: "id"},
				{Type: lex.GT, Value: ">"},
				{Type: lex.INT, Value: "10"},
			},
		},
		{
			name:  "star select",
			input: "SELECT * FROM users",
			expected: []lex.Token{
				{Type: lex.SELECT, Value: "SELECT"},
				{Type: lex.STAR, Value: "*"},
				{Type: lex.FROM, Value: "FROM"},
				{Type: lex.ID, Value: "users"},
			},
		},
		{
			name:  "special chars become IDS", // TODO: this should fail as of now the lexxer is not super strict
			input: "SELECT !!!",
			expected: []lex.Token{
				{Type: lex.SELECT, Value: "SELECT"},
				{Type: lex.ID, Value: "!!!"},
			},
		},
		{
			name:  "typo in identifer",
			input: "SELEC * FROM names",
			expected: []lex.Token{
				{Type: lex.ID, Value: "SELEC"},
				{Type: lex.STAR, Value: "*"},
				{Type: lex.FROM, Value: "FROM"},
				{Type: lex.ID, Value: "names"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lex.Lex(tt.input)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Lex(%q) = %v, expected %v", tt.input, got, tt.expected)
			}
		})
	}
}
