package app

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/ideras/exam-transcriber/transcriber"
)

type CLIOptions struct {
	OutputPath string
	PromptFile string
	BaseURL    string
	Model      string
	NoThinking bool
	ImageFiles []string
}

func parseCLIOptions(args []string, commandName string, stderr io.Writer) (CLIOptions, int) {
	flags := flag.NewFlagSet(commandName, flag.ContinueOnError)
	flags.SetOutput(stderr)

	outputPath := flags.String("output", "", "Path to the output Markdown file. If omitted, output is printed to stdout.")
	promptFile := flags.String("prompt", "", "Path to a text file containing the system prompt. (required)")
	baseURL := flags.String("base-url", "", "Base URL for an OpenAI-compatible API (optional, defaults to OpenAI)")
	model := flags.String("model", "gpt-4o", "Model name to use")
	noThinking := flags.Bool("no-thinking", false, "Disable Thinking Mode for Kimi Models (optional)")

	flags.Usage = func() {
		usage(stderr, filepath.Base(commandName))
	}

	if err := flags.Parse(args); err != nil {
		return CLIOptions{}, 2
	}

	imageFiles := flags.Args()
	if *promptFile == "" {
		fmt.Fprintln(stderr, "error: -prompt flag is required")
		flags.Usage()
		return CLIOptions{}, 1
	}

	if len(imageFiles) == 0 {
		fmt.Fprintln(stderr, "error: at least one image file must be provided")
		flags.Usage()
		return CLIOptions{}, 1
	}

	return CLIOptions{
		OutputPath: *outputPath,
		PromptFile: *promptFile,
		BaseURL:    *baseURL,
		Model:      *model,
		NoThinking: *noThinking,
		ImageFiles: imageFiles,
	}, 0
}

func usage(output io.Writer, cmd string) {
	formats := transcriber.SupportedFormats()

	fprintf := func(format string, args ...any) {
		fmt.Fprintf(output, format, args...)
	}

	fprintf("%s — Transcribe handwritten exam pages to Markdown using OpenAI Vision.\n\n", cmd)
	fprintf("USAGE:\n")
	fprintf("  %s -prompt <prompt.txt> [options] <image1> [image2 ... imageN]\n\n", cmd)

	fprintf("REQUIRED:\n")
	fprintf("  -prompt <path>         Path to system prompt text file\n")
	fprintf("  <image...>             One or more exam image files\n\n")

	fprintf("OPTIONS:\n")
	fprintf("  -output <path>         Output Markdown file (default: stdout)\n")
	fprintf("  -model <name>          Model to use (default: gpt-4o)\n")
	fprintf("  -base-url <url>        OpenAI-compatible API base URL (optional)\n")
	fprintf("  -no-thinking           Disable Thinking Mode (for Kimi models)\n\n")

	fprintf("ENVIRONMENT:\n")
	fprintf("  OPENAI_API_KEY         Required API key\n\n")

	fprintf("SUPPORTED IMAGE FORMATS:\n")
	fprintf("  %s\n\n", strings.Join(formats, ", "))

	fprintf("EXAMPLES:\n")
	fprintf("  %s -prompt docs/prompts/TranscribePrompt.v2.compact.txt page1.png\n", cmd)
	fprintf("  %s -prompt docs/prompts/TranscribePrompt.v2.txt -output exam.md page1.jpg page2.jpg\n", cmd)
	fprintf("  %s -prompt docs/prompts/TranscribePrompt.v2.code-heavy.txt -model kimi-k2 -no-thinking page1.webp\n", cmd)
	fprintf("  %s -prompt docs/prompts/TranscribePrompt.v2.txt -base-url https://api.example.com/v1 page1.png\n", cmd)
}
