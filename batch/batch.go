package batch

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"shell/executor"
)

func RunScript(scriptPath string) {
	iStream, err := os.Open(scriptPath)
	defer iStream.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(-1)
	}
	bufReader := bufio.NewReader(iStream)
	for {
		l, e := bufReader.ReadBytes('\n')
		if e == io.EOF {
			break
		}
		executor.Run(l)
	}
}
