package parser

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"shell/builtin"
	"shell/command"
	"strings"
)

// split the command into pieces of token
func split(buf []byte) []string {
	bufstr := fmt.Sprintf("%s", buf)
	bufstr = strings.TrimLeft(bufstr, " \t")
	buf = []byte(bufstr[:strings.IndexByte(bufstr, '\n')])
	if len(buf) == 0 {
		return []string{}
	}
	// replace tab with space and remove comments
	for i, c := range buf {
		if c == '\t' {
			buf[i] = ' '
		}
		if c == '#' {
			buf = buf[:i]
			break
		}
	}
	if len(buf) == 0 {
		return []string{}
	}
	bufstr = fmt.Sprintf("%s", buf)
	tokens := strings.Split(bufstr, " ")
	return tokens
}

// Parse the tokens
func ParseCMD(buf []byte) *command.CMD {
	tokens := split(buf)
	if len(tokens) == 0 {
		return nil
	}
	cmdName := tokens[0]
	cmd := command.NewCMD()
	params, in, outwr, errwr, err := parse_params(tokens[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil
	}
	if len(params) > 0 && params[len(params)-1] == "&" {
		// '&' must at the end of command if you want to set a backgroud job
		cmd.BG = true
		params = params[:len(params)-1]
	}
	if builtin.IsInternal(cmdName) {
		cmd.IsBuiltin = true
		cmd.BuiltInFunction = func() {
			builtin.BuiltInFunctionList[cmdName](
				params,
				in,
				outwr,
				errwr)
		}
	} else {
		cmd.SysCmd = exec.Command(cmdName, params...)
		cmd.SysCmd.Stdin = in
		cmd.SysCmd.Stdout = outwr
		cmd.SysCmd.Stderr = errwr
	}
	return cmd
}

func parse_params(tokens []string) ([]string,
	io.Reader,
	io.Writer,
	io.Writer,
	error) {
	in := os.Stdin
	outwr := os.Stdout
	errwr := os.Stderr
	params := []string{}
	for len(tokens) > 0 {
		switch tokens[0] {
		case "<":
			if len(tokens) == 1 {
				return params,
					in,
					outwr,
					errwr,
					fmt.Errorf("parse error near '<'")
			}
			in, _ = os.Open(tokens[1])
			if len(tokens) == 2 {
				tokens = []string{}
				break
			}
			tokens = tokens[2:]
		case ">":
			if len(tokens) == 1 {
				return params,
					in,
					outwr,
					errwr,
					fmt.Errorf("parse error near '>'")
			}
			f, _ := os.OpenFile(tokens[1],
				os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
				0664)
			outwr = f
			if len(tokens) == 2 {
				tokens = []string{}
				break
			}
			tokens = tokens[2:]
		case ">>":
			if len(tokens) == 1 {
				return params,
					in,
					outwr,
					errwr,
					fmt.Errorf("parse error near '>>'")
			}
			f, _ := os.OpenFile(tokens[1],
				os.O_WRONLY|os.O_APPEND,
				0664)
			outwr = f
			if len(tokens) == 2 {
				tokens = []string{}
				break
			}
			tokens = tokens[2:]
		default:
			param := strings.Trim(tokens[0], " ")
			if param[0] == '$' && len(param) > 1 {
				param = os.Getenv(param[1:])
			}
			if len(param) != 0 {
				params = append(params, param)
			}
			// end loop
			if len(tokens) == 1 {
				tokens = []string{}
				break
			}
			tokens = tokens[1:]
		}
	}
	return params, in, outwr, errwr, nil
}
