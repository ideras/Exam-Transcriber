package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ideras/exam-transcriber/transcriber"
	"github.com/ideras/exam-transcriber/ui"
)

func Run(args []string, commandName string) int {
	cliOptions, exitCode := parseCLIOptions(args, commandName, os.Stderr)
	if exitCode != 0 {
		return exitCode
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "error: OPENAI_API_KEY environment variable is not set")
		return 1
	}

	systemPrompt, err := transcriber.ReadPromptFile(cliOptions.PromptFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading prompt file: %v\n", err)
		return 1
	}

	fmt.Fprintf(os.Stderr, "Loading %d image(s)...\n", len(cliOptions.ImageFiles))
	contentParts, err := transcriber.BuildImageContentParts(cliOptions.ImageFiles, os.Stderr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading images: %v\n", err)
		return 1
	}

	fmt.Fprintln(os.Stderr, "Sending request to OpenAI...")
	requestTimeout := 2 * time.Minute
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	requestCtx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	start := time.Now()
	spinner := ui.StartSpinner(os.Stderr, "Transcribing exam...")
	defer spinner.Stop()

	result, err := transcriber.Transcribe(requestCtx, transcriber.Request{
		APIKey:       apiKey,
		BaseURL:      cliOptions.BaseURL,
		Model:        cliOptions.Model,
		SystemPrompt: systemPrompt,
		ContentParts: contentParts,
		NoThinking:   cliOptions.NoThinking,
	})
	elapsed := time.Since(start)

	if err != nil {
		spinner.Stop()

		switch requestCtx.Err() {
		case context.DeadlineExceeded:
			fmt.Fprintf(os.Stderr, "error calling OpenAI API: request timed out after %s\n", requestTimeout)
		case context.Canceled:
			fmt.Fprintln(os.Stderr, "error calling OpenAI API: request canceled")
		default:
			fmt.Fprintf(os.Stderr, "error calling OpenAI API: %v\n", err)
		}
		return 1
	}

	if cliOptions.OutputPath == "" {
		fmt.Fprintln(os.Stdout, result.Markdown)
	} else {
		if err := os.WriteFile(cliOptions.OutputPath, []byte(result.Markdown), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error writing output file: %v\n", err)
			return 1
		}
		fmt.Fprintf(os.Stderr, "Transcription saved to: %s\n", cliOptions.OutputPath)
	}

	fmt.Fprintf(os.Stderr, "Done. Time: %s | Tokens used — prompt: %d, completion: %d, total: %d\n",
		elapsed.Round(time.Millisecond),
		result.PromptTokens,
		result.CompletionTokens,
		result.TotalTokens,
	)

	return 0
}
