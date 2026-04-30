package exec

import (
	"fmt"
	"os"
	"strings"
)


type DropStatement struct {
	DropsWhat string
	Args []string
	FuncCall func (*DropStatement) error
}

var dropHandlers = map[string]func(*DropStatement) error {
	"DATABASE": (*DropStatement).execDropDatabase,
}

func NewDropStatement(stmt string) *DropStatement{
	stmtSplit := strings.Split(stmt, " ")
	c := &DropStatement{
		DropsWhat: stmtSplit[1],
		Args: stmtSplit[2:],
	}
	c.FuncCall = dropHandlers[strings.ToUpper(c.DropsWhat)]
	return c
}


func (c *DropStatement) execDropDatabase() error {
	name := c.Args[0]
	filename := name + ".db"
	exists := doesDatabaseNameExist(name)
	if !exists {
		return fmt.Errorf("Database with name %s not found", name)
	}
	err := os.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}
