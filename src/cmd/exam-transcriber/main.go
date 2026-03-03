package main

import (
	"os"
	"path/filepath"

	"github.com/ideras/exam-transcriber/app"
)

func main() {
	exitCode := app.Run(os.Args[1:], filepath.Base(os.Args[0]))
	os.Exit(exitCode)
}
