package parser

import (
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
			err := checkCharacters(tt.tokens)
			gotError := err != nil
			if gotError != tt.expectedError {
				t.Errorf("CheckCharacters(%v) error = %v, expectedError %v", tt.tokens, gotError, tt.expectedError)
			}
		})
	}

}
