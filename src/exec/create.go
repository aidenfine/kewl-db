package exec

import (
	"fmt"
	"strings"

	"github.com/aidenfine/kewl-db/src/btree"
)

type CreateStatement struct {
	CreatesWhat string
	Args        []string
	FuncCall    func(*CreateStatement) error
}

var createHandlers = map[string]func(*CreateStatement) error{
	"DATABASE": (*CreateStatement).execCreateDatabase,
}

func NewCreateStatement(stmt string) *CreateStatement {
	stmtSplit := strings.Split(stmt, " ")
	c := &CreateStatement{
		CreatesWhat: stmtSplit[1],
		Args:        stmtSplit[2:],
	}
	c.FuncCall = createHandlers[strings.ToUpper(c.CreatesWhat)]
	return c
}

func (c *CreateStatement) Exec() error {
	switch strings.ToUpper(c.CreatesWhat) {
	case "DATABASE":
		return c.execCreateDatabase()
	default:
		return fmt.Errorf("unknown CREATE target: %s", c.CreatesWhat)
	}
}

func (c *CreateStatement) execCreateDatabase() error {
	name := c.Args[0]
	exists := doesDatabaseNameExist(name)
	if exists {
		return fmt.Errorf("Database with name %s already exists", name)
	}
	_, err := btree.NewBTree(2, name+".db")
	fmt.Println("Created database")
	return err
}
