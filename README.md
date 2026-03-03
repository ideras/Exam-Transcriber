# exam-transcriber — Handwritten Exam OCR to Markdown (OpenAI Vision CLI in Go)

![Go Version](https://img.shields.io/badge/Go-1.22+-blue)
![OpenAI Compatible](https://img.shields.io/badge/OpenAI-Compatible-purple)
![CLI Tool](https://img.shields.io/badge/type-CLI-orange)
![Academic Use](https://img.shields.io/badge/use-academic%20workflows-green)

**exam-transcriber** is a Go-based CLI tool that converts handwritten exam scans into structured Markdown using OpenAI-compatible multimodal (vision) models such as `gpt-4o`.

It is designed for academic digitization workflows, AI-assisted grading pipelines, and large-scale transcription of handwritten student submissions.

---

## ✨ Features

- Handwritten exam transcription to Markdown
- Vision-model powered OCR alternative
- Multi-page processing in a single invocation
- Model-agnostic (any OpenAI-compatible endpoint)
- Prompt-driven behavior (course-specific customization)
- Clean stdout output (pipe-friendly)
- Token usage reporting to stderr
- Batch processing support

---

## Why This Exists

In the era of generative AI, computer-based programming exams have become increasingly difficult to administer in ways that ensure independent student work.

Handwritten, paper-based assessments remain an effective way to evaluate reasoning and authorship. However, grading paper exams—especially those containing code, algorithms, and structured solutions—is operationally intensive.

`exam-transcriber` addresses this gap by converting handwritten submissions into structured Markdown, enabling scalable and controlled AI-assisted evaluation workflows while preserving academic integrity.

---

## Comparison to Traditional OCR

| Feature | Traditional OCR (e.g., Tesseract) | exam-transcriber |
|----------|-----------------------------------|------------------|
| Handwriting handling | Limited | Strong (LLM-based) |
| Code block recovery | No | Yes |
| Math formatting | Raw text | Structured Markdown |
| Context awareness | No | Yes |
| Prompt control | No | Yes |

Unlike conventional OCR engines, exam-transcriber leverages multimodal LLMs to produce semantically structured output suitable for grading pipelines.

---

## Requirements

- Go 1.22+
- OpenAI-compatible API key
- Access to a vision-capable model (default: `gpt-4o`)

---

## Installation

```bash
git clone <your-repo>
cd exam-transcriber
go build -o exam-transcriber ./cmd/exam-transcriber
````

Optionally move the binary to your `$PATH`.

---

## Configuration

Set your API key:

```bash
export OPENAI_API_KEY="sk-..."
```

If using a compatible provider:

```bash
export OPENAI_BASE_URL="https://api.example.com/v1"
```

---

## Usage

```bash
exam-transcriber -prompt <prompt-file> [options] <image1> [image2 ... imageN]
```

Images are processed in the order provided.

---

## Flags

| Flag           | Required | Description                                           |
| -------------- | -------- | ----------------------------------------------------- |
| `-prompt`      | Yes      | Path to the system prompt file                        |
| `-output`      | No       | Output Markdown file (defaults to stdout)             |
| `-model`       | No       | Model name (default: `gpt-4o`)                        |
| `-base-url`    | No       | OpenAI-compatible API endpoint                        |
| `-no-thinking` | No       | Disable reasoning mode (useful for certain providers) |

---

## Examples

### Transcribe a Single Page

```bash
./exam-transcriber -prompt docs/prompts/prompt.txt exam_page1.png
```

### Transcribe Multiple Pages

```bash
./exam-transcriber \
  -prompt docs/prompts/prompt.txt \
  -output student_john.md \
  page1.jpg page2.jpg page3.jpg
```

### Batch Process Multiple Students

```bash
for student in exams/*/; do
  name=$(basename "$student")
  ./exam-transcriber \
    -prompt docs/prompts/prompt.txt \
    -output "output/${name}.md" \
    "$student"*.jpg
done
```

### Use a Custom Model

```bash
./exam-transcriber \
  -prompt docs/prompts/prompt.txt \
  -model kimi-k2 \
  page1.png
```

### Use an OpenAI-Compatible Endpoint

```bash
./exam-transcriber \
  -prompt docs/prompts/prompt.txt \
  -base-url https://api.example.com/v1 \
  page1.png
```

---

## Supported Image Formats

* `jpg`
* `jpeg`
* `png`
* `gif`
* `webp`

**Recommended:** Scan at 300 DPI or higher with minimal skew for best results.

---

## Prompt Customization

The `-prompt` file is sent as the system instruction before images are processed.

This enables:

* Course-specific formatting rules
* Rubric-aware transcription
* Code block enforcement
* Markdown normalization
* Strict output constraints

Prompt variants are located in:

```
docs/prompts/
```

This design keeps evaluation policy separate from application logic.

---

## Output Behavior

* Markdown output → stdout (default) or `-output` file
* Logs and token usage → stderr

This ensures clean piping:

```bash
./exam-transcriber ... > submission.md
```

Example stderr output:

```
Loading 4 image(s)...
Sending request to OpenAI...
Transcription saved to: student_john.md
Done. Tokens used — prompt: 8432, completion: 1205, total: 9637
```

---

## Cost Estimate

Using `gpt-4o` with high-detail image processing:

Typical 4–6 page handwritten exam:
**~$0.05 – $0.15 USD**

Actual cost depends on:

* Image resolution
* Answer length
* Selected model
* Provider pricing

---

## Architecture Overview

The project follows a clean layered structure:

```
cmd/exam-transcriber/main.go   # CLI entrypoint
internal/app/                  # flag parsing and orchestration
internal/transcriber/          # prompt + image loading + API request
internal/ui/                   # terminal spinner
scripts/build_release.sh       # release build script
docs/prompts/                  # system prompt variants
```

The architecture allows easy model swapping and prompt experimentation.

---

## Use Cases

* University exam digitization
* AI-assisted grading workflows
* TA evaluation automation
* Converting scanned homework to Markdown
* Code-heavy exam transcription
* Multimodal OCR alternative for academic use

---

## License

MIT
