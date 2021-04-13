package builtin

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type BuiltInFunc func([]string, io.Reader, io.Writer, io.Writer)

var BuiltInFunctionList map[string]BuiltInFunc = make(map[string]BuiltInFunc)

func IsInternal(cmd string) bool {
	_, ok := BuiltInFunctionList[cmd]
	return ok
}

func init() {
	BuiltInFunctionList["cd"] = cd
	BuiltInFunctionList["exit"] = exit
	BuiltInFunctionList["export"] = export
}

func cd(args []string, in io.Reader, outwr io.Writer, errwr io.Writer) {
	if len(args) == 0 {
		args = append(args, os.Getenv("HOME"))
	}
	if len(args) > 1 {
		errwr.Write([]byte("too much arguments\n"))
		return
	}
	err := os.Chdir(args[0])
	pushDIR(args[0])
	if err != nil {
		fmt.Fprintf(errwr, "%s\n", err)
	}
}

func exit(args []string, in io.Reader, outwr io.Writer, errwr io.Writer) {
	os.Exit(0)
}

func export(args []string, in io.Reader, outwr io.Writer, errwr io.Writer) {
	for _, str := range args {
		tokens := strings.Split(str, "=")
		if len(tokens) != 2 {
			errwr.Write([]byte("An error occurred!\n"))
		}
		os.Setenv(tokens[0], tokens[1])
	}
}
