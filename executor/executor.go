package executor

import (
	"fmt"
	"os/exec"
	"shell/parser"
)

func Run(cmdBytes []byte) {
	cmd := parser.ParseCMD(cmdBytes)
	if cmd == nil {
		return
	}
	if cmd.IsBuiltin {
		cmd.BuiltInFunction()
	} else if cmd.BG {
		cmd.SysCmd.Start()
	} else {
		err := cmd.SysCmd.Run()
		if err != nil {
			switch err.(type) {
			case *exec.Error:
				fmt.Fprintf(cmd.SysCmd.Stderr, "gosh: command not found: %s\n", cmd.SysCmd.Path)
			default:
				fmt.Fprintf(cmd.SysCmd.Stderr, "%s\n", err)
			}
		}
	}
}
