package command

import (
	"os/exec"
)

type CMD struct {
	SysCmd *exec.Cmd

	IsBuiltin       bool
	BuiltInFunction func()
	BG              bool
}

func NewCMD() *CMD {
	return &CMD{&exec.Cmd{}, false, nil, false}
}
