package exec
import (
	"errors"
	"os"
)


// todo refactor this to be a bit safer in case there is an error
func doesDatabaseNameExist(name string) bool {
	_, err := os.Stat(name+".db")
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist){
		return false
	}
	return false
}
