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

func TestCheckCharacters(t *testing.T) {
	tests := []struct {
		name          string
		tokens        []lex.Token
		expectedError bool
	}{
		{
			name:          "letters and numbers",
			tokens:        []lex.Token{{Value: "users42"}},
			expectedError: false,
		},
		{
			name:          "hyphens and underscores",
			tokens:        []lex.Token{{Value: "user-name_2"}},
			expectedError: false,
		},
		{
			name:          "unicode letters and numbers",
			tokens:        []lex.Token{{Value: "café２０２４"}},
			expectedError: false,
		},
		{
			name:          "empty token list",
			tokens:        nil,
			expectedError: false,
		},
		{
			name: "multiple valid tokens",
			tokens: []lex.Token{
				{Value: "SELECT"},
				{Value: "account_id"},
				{Value: "10"},
			},
			expectedError: false,
		},
		{
			name:          "punctuation",
			tokens:        []lex.Token{{Value: "user.name"}},
			expectedError: true,
		},
		{
			name:          "whitespace",
			tokens:        []lex.Token{{Value: "user name"}},
			expectedError: true,
		},
		{
			name:          "invalid UTF-8",
			tokens:        []lex.Token{{Value: string([]byte{0xff})}},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := lex.CheckCharacters(tt.tokens)
			gotError := err != nil
			if gotError != tt.expectedError {
				t.Errorf("CheckCharacters(%v) error = %v, expectedError %v", tt.tokens, gotError, tt.expectedError)
			}
		})
	}

}
