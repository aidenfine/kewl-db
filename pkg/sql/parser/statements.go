package parser

import (
	"fmt"

	"github.com/aidenfine/kewl-db/pkg/sql/ast"
	"github.com/aidenfine/kewl-db/pkg/sql/lex"
)

// Statement-specific parsers are called after dispatch has checked the first token.
func (p *Parser) parseInsert() (*ast.InsertStatement, error) {
	p.pos++ // INSERT
	if _, err := p.expect(lex.INTO, "INTO"); err != nil {
		return nil, err
	}
	table, err := p.expect(lex.ID, "table name")
	if err != nil {
		return nil, err
	}
	stmt := &ast.InsertStatement{Table: table.Value}
	if p.match(lex.LPAREN) {
		for {
			column, err := p.expect(lex.ID, "column name")
			if err != nil {
				return nil, err
			}
			stmt.Columns = append(stmt.Columns, column.Value)
			if !p.match(lex.COMMA) {
				break
			}
		}
		if _, err := p.expect(lex.RPAREN, ")"); err != nil {
			return nil, err
		}
	}
	if _, err := p.expect(lex.VALUES, "VALUES"); err != nil {
		return nil, err
	}
	width := len(stmt.Columns)
	for {
		if _, err := p.expect(lex.LPAREN, "("); err != nil {
			return nil, err
		}
		var row []ast.Expr
		for {
			value, err := p.parseLiteral()
			if err != nil {
				return nil, err
			}
			row = append(row, value)
			if !p.match(lex.COMMA) {
				break
			}
		}
		if _, err := p.expect(lex.RPAREN, ")"); err != nil {
			return nil, err
		}
		if width == 0 {
			width = len(row)
		}
		if len(row) != width {
			return nil, fmt.Errorf("%w: VALUES row %d has %d values, expected %d", ErrSyntaxError, len(stmt.Values)+1, len(row), width)
		}
		stmt.Values = append(stmt.Values, row)
		if !p.match(lex.COMMA) {
			break
		}
	}
	return stmt, nil
}

func (p *Parser) parseUpdate() (*ast.UpdateStatement, error) {
	p.pos++ // UPDATE
	table, err := p.expect(lex.ID, "table name")
	if err != nil {
		return nil, err
	}
	if _, err := p.expect(lex.SET, "SET"); err != nil {
		return nil, err
	}
	stmt := &ast.UpdateStatement{Table: table.Value}
	for {
		column, err := p.expect(lex.ID, "column name")
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(lex.EQ, "="); err != nil {
			return nil, err
		}
		value, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		stmt.Assignments = append(stmt.Assignments, ast.Assignment{Column: column.Value, Value: value})
		if !p.match(lex.COMMA) {
			break
		}
	}
	stmt.Where, err = p.parseOptionalWhere()
	if err != nil {
		return nil, err
	}
	return stmt, nil
}

func (p *Parser) parseDelete() (*ast.DeleteStatement, error) {
	p.pos++ // DELETE
	if _, err := p.expect(lex.FROM, "FROM"); err != nil {
		return nil, err
	}
	table, err := p.expect(lex.ID, "table name")
	if err != nil {
		return nil, err
	}
	where, err := p.parseOptionalWhere()
	if err != nil {
		return nil, err
	}
	return &ast.DeleteStatement{Table: table.Value, Where: where}, nil
}

func (p *Parser) parseOptionalWhere() (ast.Expr, error) {
	if p.match(lex.WHERE) {
		return p.parseComparison()
	}
	return nil, nil
}

func (p *Parser) parseDatabase() (ast.Statement, error) {
	kind := p.tokens[p.pos].Type
	p.pos++ // CREATE or DROP
	if _, err := p.expect(lex.DATABASE, "DATABASE"); err != nil {
		return nil, err
	}
	name, err := p.expect(lex.ID, "database name")
	if err != nil {
		return nil, err
	}
	if kind == lex.CREATE {
		return &ast.CreateDatabaseStatement{Name: name.Value}, nil
	}
	return &ast.DropDatabaseStatement{Name: name.Value}, nil
}
