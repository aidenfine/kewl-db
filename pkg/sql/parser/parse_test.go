package parser

import (
	"reflect"
	"testing"

	"github.com/aidenfine/kewl-db/pkg/sql/lex"
)

func TestCheckCharacters(t *testing.T) {
	tests := []struct {
		name          string
		tokens        []lex.Token
		expectedError bool
	}{
		{
			name:          "letters and numbers",
			tokens:        []lex.Token{{Type: lex.ID, Value: "users42"}},
			expectedError: false,
		},
		{
			name:          "hyphens and underscores",
			tokens:        []lex.Token{{Type: lex.ID, Value: "user-name_2"}},
			expectedError: false,
		},
		{
			name:          "unicode letters and numbers",
			tokens:        []lex.Token{{Type: lex.ID, Value: "café２０２４"}},
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
				{Type: lex.SELECT, Value: "SELECT"},
				{Type: lex.ID, Value: "account_id"},
				{Type: lex.INT, Value: "10"},
			},
			expectedError: false,
		},
		{
			name:          "punctuation",
			tokens:        []lex.Token{{Type: lex.ID, Value: "user.name"}},
			expectedError: true,
		},
		{
			name:          "whitespace",
			tokens:        []lex.Token{{Type: lex.ID, Value: "user name"}},
			expectedError: true,
		},
		{
			name:          "invalid UTF-8",
			tokens:        []lex.Token{{Type: lex.ID, Value: string([]byte{0xff})}},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkCharacters(tt.tokens)
			gotError := err != nil
			if gotError != tt.expectedError {
				t.Errorf("CheckCharacters(%v) error = %v, expectedError %v", tt.tokens, gotError, tt.expectedError)
			}
		})
	}

}

func TestStripAndCheckParens(t *testing.T) {
	tests := []struct {
		name           string
		tokens         []lex.Token
		expectedTokens []lex.Token
		expectedError  bool
	}{
		{
			name: "no parentheses",
			tokens: []lex.Token{
				{Type: lex.SELECT, Value: "SELECT"},
				{Type: lex.ID, Value: "account_id"},
				{Type: lex.INT, Value: "10"},
			},
			expectedTokens: []lex.Token{
				{Type: lex.SELECT, Value: "SELECT"},
				{Type: lex.ID, Value: "account_id"},
				{Type: lex.INT, Value: "10"},
			},
			expectedError: false,
		},
		{
			name: "balanced parentheses",
			tokens: []lex.Token{
				{Value: "("},
				{Type: lex.ID, Value: "account_id"},
				{Value: ")"},
			},
			expectedTokens: []lex.Token{{Type: lex.ID, Value: "account_id"}},
			expectedError:  false,
		},
		{
			name: "nested mixed parentheses",
			tokens: []lex.Token{
				{Value: "("}, {Value: "["},
				{Type: lex.ID, Value: "value"},
				{Value: "]"}, {Value: ")"},
			},
			expectedTokens: []lex.Token{{Type: lex.ID, Value: "value"}},
			expectedError:  false,
		},
		{
			name:          "mismatched parentheses",
			tokens:        []lex.Token{{Value: "("}, {Value: "]"}},
			expectedError: true,
		},
		{
			name:          "closing parenthesis without opener",
			tokens:        []lex.Token{{Value: ")"}},
			expectedError: true,
		},
		{
			name:          "unclosed parenthesis",
			tokens:        []lex.Token{{Value: "("}},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, gotTokens, err := stripAndCheckParens(tt.tokens)
			gotError := err != nil

			if gotError != tt.expectedError {
				t.Errorf("stripAndCheckParens(%v) error = %v, expectedError %v", tt.tokens, gotError, tt.expectedError)
			}
			if !tt.expectedError {
				if !valid {
					t.Errorf("stripAndCheckParens(%v) valid = false, expected true", tt.tokens)
				}
				if !reflect.DeepEqual(gotTokens, tt.expectedTokens) {
					t.Errorf("stripAndCheckParens(%v) tokens = %v, expected %v", tt.tokens, gotTokens, tt.expectedTokens)
				}
			}
		})
	}
}
