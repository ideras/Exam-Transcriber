# exam-transcriber

A CLI tool that transcribes handwritten student exam pages into Markdown using an OpenAI-compatible vision model (default: `gpt-4o`).

## Prerequisites

- Go 1.22+
- An OpenAI API key with access to your selected model (default: `gpt-4o`)

## Installation

```bash
git clone <your-repo>
cd exam-transcriber
go build -o exam-transcriber ./cmd/exam-transcriber
```

## Configuration

Set your OpenAI API key as an environment variable:

```bash
export OPENAI_API_KEY="sk-..."
```

## Usage

```
exam-transcriber -prompt <docs/prompts/prompt.txt> [options] <image1> [image2 ... imageN]
```

### Flags

| Flag           | Required | Description                                                      |
|----------------|----------|------------------------------------------------------------------|
| `-prompt`      | Yes      | Path to a text file containing the system prompt                 |
| `-output`      | No       | Path to the output Markdown file (defaults to stdout)            |
| `-model`       | No       | Model name to use (default: `gpt-4o`)                            |
| `-base-url`    | No       | OpenAI-compatible API base URL                                   |
| `-no-thinking` | No       | Disable Thinking Mode (for Kimi models)                          |

### Examples

Transcribe a single page and print to stdout:
```bash
./exam-transcriber -prompt docs/prompts/prompt.txt exam_page1.png
```

Transcribe a multi-page exam and save to a file:
```bash
./exam-transcriber -prompt docs/prompts/prompt.txt -output student_john.md \
  page1.jpg page2.jpg page3.jpg page4.jpg
```

Batch-transcribe multiple students using a shell loop:
```bash
for student in exams/*/; do
  name=$(basename "$student")
  ./exam-transcriber -prompt docs/prompts/prompt.txt \
    -output "output/${name}.md" \
    "$student"*.jpg
done
```

Use a custom model:
```bash
./exam-transcriber -prompt docs/prompts/prompt.txt -model kimi-k2 page1.png
```

Use an OpenAI-compatible endpoint:
```bash
./exam-transcriber -prompt docs/prompts/prompt.txt -base-url https://api.example.com/v1 page1.png
```

Using curated prompt variants:
```bash
./exam-transcriber -prompt docs/prompts/TranscribePrompt.v2.compact.txt page1.png
./exam-transcriber -prompt docs/prompts/TranscribePrompt.v2.code-heavy.txt page1.png
```

## Supported Image Formats

`jpg`, `jpeg`, `png`, `gif`, `webp`

> **Tip:** Scan exams at 300 DPI or higher for best transcription accuracy.

## Prompt Customization

The `-prompt` file is the system prompt sent to the selected model before the images.
A ready-to-use example is provided in `docs/prompts/prompt.txt`. You can create different
prompt files for different courses or exam formats without touching the code.

Available prompt variants are under `docs/prompts/`.

## Project Layout

```text
cmd/exam-transcriber/main.go   # CLI entrypoint
internal/app/                  # flag parsing + orchestration
internal/transcriber/          # prompt/image loading + OpenAI request
internal/ui/                   # terminal spinner
scripts/build_release.sh       # release build script
docs/prompts/                  # system prompt variants
```

## Output

The tool writes the raw Markdown response from the selected model to the specified output
file (or stdout). Progress and token usage are printed to stderr so they don't
interfere with piped output.

Example stderr output:
```
Loading 4 image(s)...
  Loaded: page1.jpg (image/jpeg, 412 KB)
  Loaded: page2.jpg (image/jpeg, 389 KB)
  Loaded: page3.jpg (image/jpeg, 401 KB)
  Loaded: page4.jpg (image/jpeg, 375 KB)
Sending request to OpenAI...
Transcription saved to: student_john.md
Done. Tokens used — prompt: 8432, completion: 1205, total: 9637
```

## Cost Estimate

Using `gpt-4o` with `detail: high` per image, a typical 4–6 page exam costs
roughly **$0.05–$0.15** depending on image resolution and answer length.

Costs vary by model and provider if you set `-model` or `-base-url`.
