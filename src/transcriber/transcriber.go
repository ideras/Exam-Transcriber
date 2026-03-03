package transcriber

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type Request struct {
	APIKey       string
	BaseURL      string
	Model        string
	SystemPrompt string
	ContentParts []openai.ChatCompletionContentPartUnionParam
	NoThinking   bool
}

type Result struct {
	Markdown         string
	PromptTokens     int64
	CompletionTokens int64
	TotalTokens      int64
}

func ReadPromptFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func BuildImageContentParts(paths []string, logOutput io.Writer) ([]openai.ChatCompletionContentPartUnionParam, error) {
	parts := make([]openai.ChatCompletionContentPartUnionParam, 0, len(paths))

	for _, path := range paths {
		encodedImage, err := encodeImageDataURL(path)
		if err != nil {
			return nil, err
		}

		parts = append(parts, openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
			URL:    encodedImage.DataURL,
			Detail: "high",
		}))

		fmt.Fprintf(logOutput, "  Loaded: %s (%s, %d KB)\n", filepath.Base(path), encodedImage.MIMEType, encodedImage.SizeKB)
	}

	return parts, nil
}

func Transcribe(ctx context.Context, request Request) (Result, error) {
	clientOptions := []option.RequestOption{option.WithAPIKey(request.APIKey)}
	if strings.TrimSpace(request.BaseURL) != "" {
		clientOptions = append(clientOptions, option.WithBaseURL(request.BaseURL))
	}
	client := openai.NewClient(clientOptions...)

	userPrompt := "Please transcribe all pages of this exam as instructed."
	textPart := openai.TextContentPart(userPrompt)
	allParts := append([]openai.ChatCompletionContentPartUnionParam{textPart}, request.ContentParts...)

	params := openai.ChatCompletionNewParams{
		Model: request.Model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(request.SystemPrompt),
			{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfArrayOfContentParts: allParts,
					},
				},
			},
		},
		Temperature: openai.Float(1.0),
	}

	if request.NoThinking {
		params.SetExtraFields(map[string]any{
			"thinking": map[string]any{
				"type": "disabled",
			},
		})
	}

	response, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		return Result{}, err
	}

	if len(response.Choices) == 0 {
		return Result{}, errors.New("no response choices returned from OpenAI")
	}

	return Result{
		Markdown:         response.Choices[0].Message.Content,
		PromptTokens:     int64(response.Usage.PromptTokens),
		CompletionTokens: int64(response.Usage.CompletionTokens),
		TotalTokens:      int64(response.Usage.TotalTokens),
	}, nil
}
