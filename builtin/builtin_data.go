package builtin

import (
	"fmt"
	"os"
)

var dirStack []string

func init() {
	dirStack = append(dirStack, os.Getenv("PWD"))
}

func pushDIR(dir string) {
	dirStack = append(dirStack, dir)
}

// dirStack[0] is reserved to save the directory where the shell is started
func popDIR() (string, error) {
	if len(dirStack) == 1 {
		return "", fmt.Errorf("directory stack empty")
	}
	dir := dirStack[len(dirStack)-1]
	dirStack = dirStack[0 : len(dirStack)-1]
	return dir, nil
}
func topDIR() string {
	return dirStack[len(dirStack)-1]
}
