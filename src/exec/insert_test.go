package exec

import (
	"testing"
)

func TestNewInsertStatement(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		table   string
		columns []string
		values  []string
		wantErr bool
	}{
		{
			name:    "full insert with columns",
			input:   "INSERT INTO users (name, age) VALUES (alice, 30)",
			table:   "users",
			columns: []string{"name", "age"},
			values:  []string{"alice", "30"},
		},
		{
			name:   "insert without columns",
			input:  "INSERT INTO users VALUES (alice, 30)",
			table:  "users",
			columns: nil,
			values: []string{"alice", "30"},
		},
		{
			name:    "single value",
			input:   "INSERT INTO config (key) VALUES (debug)",
			table:   "config",
			columns: []string{"key"},
			values:  []string{"debug"},
		},
		{
			name:    "missing VALUES keyword",
			input:   "INSERT INTO users (name) (alice)",
			wantErr: true,
		},
		{
			name:    "missing INTO keyword",
			input:   "INSERT users VALUES (alice)",
			wantErr: true,
		},
		{
			name:    "column count mismatch",
			input:   "INSERT INTO users (name, age) VALUES (alice)",
			wantErr: true,
		},
		{
			name:    "missing table name",
			input:   "INSERT INTO  VALUES (alice)",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ins, err := NewInsertStatement(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ins.Table != tt.table {
				t.Errorf("table = %q, want %q", ins.Table, tt.table)
			}
			if len(ins.Columns) != len(tt.columns) {
				t.Fatalf("columns = %v, want %v", ins.Columns, tt.columns)
			}
			for i, c := range ins.Columns {
				if c != tt.columns[i] {
					t.Errorf("columns[%d] = %q, want %q", i, c, tt.columns[i])
				}
			}
			if len(ins.Values) != len(tt.values) {
				t.Fatalf("values = %v, want %v", ins.Values, tt.values)
			}
			for i, v := range ins.Values {
				if v != tt.values[i] {
					t.Errorf("values[%d] = %q, want %q", i, v, tt.values[i])
				}
			}
		})
	}
}
