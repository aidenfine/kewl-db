package main

import (
	"errors"
	"os"
)

// SELECT col1, col2,col3 FROM table
func Select(cols []string, table string) {
	// TOOD
}

func CreateDBFile(name string) error {
	for _, c := range name {
		if !IsAlphaNumeric(byte(c)) {
			return errors.New("Non alnum char found")
		}
	}

	filename := name + ".db"
	err := os.WriteFile(filename, []byte(""), 0755)
	if err != nil {
		return err
	}
	return nil
}
