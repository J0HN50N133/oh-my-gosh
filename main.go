package main

import (
	"fmt"
	"io"
	"os"
	"shell/batch"
	"shell/executor"
)

var (
	PROMPT     string = "\033[1;31;31m>>>\033[0m"
	MAX_BUFFER int    = 512
)

func main() {
	argc := len(os.Args)
	if argc > 2 {
		fmt.Fprintf(os.Stderr, "An error has occurred!")
		os.Exit(-1)
	} else if argc > 1 {
		batch.RunScript(os.Args[1])
	} else {
		for {
			showPrompt(PROMPT)
			buf := make([]byte, MAX_BUFFER)
			iStream := os.Stdin
			_, err := iStream.Read(buf)
			if err == io.EOF {
				os.Exit(0)
			}
			// cmd := parser.ParseCMD(buf)
			executor.Run(buf)
		}
	}
}

func showPrompt(prompt string) {
	fmt.Printf("%s ", prompt)
}
