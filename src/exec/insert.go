package exec

import (
	"fmt"
	"strings"
)

type InsertStatement struct {
	Table   string
	Columns []string
	Values  []string
}

func NewInsertStatement(stmt string) (*InsertStatement, error) {
	upper := strings.ToUpper(stmt)

	if !strings.HasPrefix(upper, "INSERT INTO ") {
		return nil, fmt.Errorf("syntax error: expected INSERT INTO")
	}

	rest := stmt[len("INSERT INTO "):]

	valuesIdx := strings.Index(strings.ToUpper(rest), "VALUES")
	if valuesIdx == -1 {
		return nil, fmt.Errorf("syntax error: expected VALUES")
	}

	tablePart := strings.TrimSpace(rest[:valuesIdx])
	valuesPart := strings.TrimSpace(rest[valuesIdx+len("VALUES"):])

	table, columns, err := parseTableAndColumns(tablePart)
	if err != nil {
		return nil, err
	}

	values, err := parseParenList(valuesPart)
	if err != nil {
		return nil, fmt.Errorf("syntax error in VALUES: %w", err)
	}

	if len(columns) > 0 && len(columns) != len(values) {
		return nil, fmt.Errorf("column count (%d) does not match value count (%d)", len(columns), len(values))
	}

	return &InsertStatement{
		Table:   table,
		Columns: columns,
		Values:  values,
	}, nil
}

func parseTableAndColumns(s string) (string, []string, error) {
	parenIdx := strings.Index(s, "(")
	if parenIdx == -1 {
		table := strings.TrimSpace(s)
		if table == "" {
			return "", nil, fmt.Errorf("syntax error: missing table name")
		}
		return table, nil, nil
	}

	table := strings.TrimSpace(s[:parenIdx])
	if table == "" {
		return "", nil, fmt.Errorf("syntax error: missing table name")
	}

	columns, err := parseParenList(s[parenIdx:])
	if err != nil {
		return "", nil, fmt.Errorf("syntax error in column list: %w", err)
	}

	return table, columns, nil
}

func parseParenList(s string) ([]string, error) {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "(") || !strings.HasSuffix(s, ")") {
		return nil, fmt.Errorf("expected parenthesized list, got: %s", s)
	}

	inner := s[1 : len(s)-1]
	parts := strings.Split(inner, ",")

	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed == "" {
			return nil, fmt.Errorf("empty item in list")
		}
		result = append(result, trimmed)
	}

	return result, nil
}
