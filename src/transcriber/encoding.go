package transcriber

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var supportedExtensions = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".webp": "image/webp",
}

type EncodedImage struct {
	DataURL  string
	MIMEType string
	SizeKB   int
}

func SupportedFormats() []string {
	formats := make([]string, 0, len(supportedExtensions))
	for extension := range supportedExtensions {
		formats = append(formats, strings.TrimPrefix(extension, "."))
	}
	sort.Strings(formats)
	return formats
}

func encodeImageDataURL(path string) (EncodedImage, error) {
	extension := strings.ToLower(filepath.Ext(path))
	mimeType, ok := supportedExtensions[extension]
	if !ok {
		return EncodedImage{}, fmt.Errorf("unsupported extension %q", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return EncodedImage{}, fmt.Errorf("cannot read file %q: %w", path, err)
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	return EncodedImage{
		DataURL:  fmt.Sprintf("data:%s;base64,%s", mimeType, encoded),
		MIMEType: mimeType,
		SizeKB:   len(data) / 1024,
	}, nil
}
