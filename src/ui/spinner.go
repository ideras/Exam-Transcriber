package ui

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type Spinner struct {
	done chan struct{}
	once sync.Once
}

func StartSpinner(output io.Writer, label string) *Spinner {
	spinner := &Spinner{done: make(chan struct{})}

	go func() {
		frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		frameIndex := 0
		for {
			select {
			case <-spinner.done:
				fmt.Fprintf(output, "\r✓ %s\n", label)
				return
			default:
				fmt.Fprintf(output, "\r%s %s", frames[frameIndex%len(frames)], label)
				time.Sleep(80 * time.Millisecond)
				frameIndex++
			}
		}
	}()

	return spinner
}

func (spinner *Spinner) Stop() {
	spinner.once.Do(func() {
		close(spinner.done)
	})
}
