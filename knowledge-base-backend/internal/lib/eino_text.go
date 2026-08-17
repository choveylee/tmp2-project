package lib

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/choveylee/terror"
	"github.com/choveylee/tlog"
	markdownSplitter "github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/schema"

	constant "dev.choveylee.top/knowledge-base-backend/internal/const"
)

func SplitEinoText(ctx context.Context, documentId string, plainText string) ([]*schema.Document, *terror.Terror) {
	normalizedText := NormalizeDocumentText(plainText)
	sourceText := normalizedText
	if !hasMarkdownHeading(normalizedText) {
		sourceText = buildMarkdownSections(normalizedText, einoTextChunkSize)
	}

	sourceDocument := &schema.Document{
		ID:      documentId,
		Content: sourceText,
		MetaData: map[string]any{
			"document_id": documentId,
		},
	}

	splitter, err := markdownSplitter.NewHeaderSplitter(ctx, &markdownSplitter.HeaderConfig{
		Headers: map[string]string{
			"#":   "h1",
			"##":  "h2",
			"###": "h3",
		},
		TrimHeaders: true,
		IDGenerator: func(ctx context.Context, originalID string, splitIndex int) string {
			return fmt.Sprintf("%s_chunk_%03d", originalID, splitIndex+1)
		},
	})
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Split eino text (document id: %s, text len: %d, chunk size: %d) err (new header splitter %v)",
			documentId, utf8.RuneCountInString(normalizedText), einoTextChunkSize, err)

		errx := terror.NewTerror(ctx, err, constant.ErrorCodeDocumentParseFailed, errMsg)

		return nil, errx
	}

	chunks, err := splitter.Transform(ctx, []*schema.Document{sourceDocument})
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Split eino text (document id: %s, text len: %d, chunk size: %d) err (splitter transform %v)",
			documentId, utf8.RuneCountInString(normalizedText), einoTextChunkSize, err)

		errx := terror.NewTerror(ctx, err, constant.ErrorCodeDocumentParseFailed, errMsg)

		return nil, errx
	}

	filteredChunks := make([]*schema.Document, 0, len(chunks))
	for _, chunk := range chunks {
		chunk.Content = NormalizeDocumentText(chunk.Content)
		if chunk.Content == "" {
			continue
		}

		filteredChunks = append(filteredChunks, chunk)
	}

	return filteredChunks, nil
}

func NormalizeDocumentText(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	lines := strings.Split(text, "\n")
	normalizedLines := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		normalizedLines = append(normalizedLines, line)
	}

	return strings.TrimSpace(strings.Join(normalizedLines, "\n"))
}

func EstimateDocumentTokenCount(text string) uint {
	return uint(utf8.RuneCountInString(text))
}

func Sha256Text(text string) string {
	hash := sha256.Sum256([]byte(text))

	return hex.EncodeToString(hash[:])
}

func hasMarkdownHeading(text string) bool {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") || strings.HasPrefix(line, "## ") || strings.HasPrefix(line, "### ") {
			return true
		}
	}

	return false
}

func buildMarkdownSections(text string, chunkSize int) string {
	if strings.TrimSpace(text) == "" {
		return ""
	}

	paragraphs := strings.Split(text, "\n")
	sections := make([]string, 0)
	currentParagraphs := make([]string, 0)
	currentLen := 0
	sectionNo := 1

	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if paragraph == "" {
			continue
		}

		paragraphLen := utf8.RuneCountInString(paragraph)
		if currentLen > 0 && currentLen+paragraphLen > chunkSize {
			sections = append(sections, fmt.Sprintf("## Part %03d\n%s", sectionNo, strings.Join(currentParagraphs, "\n")))
			sectionNo++

			currentParagraphs = currentParagraphs[:0]
			currentLen = 0
		}

		currentParagraphs = append(currentParagraphs, paragraph)
		currentLen += paragraphLen
	}

	if len(currentParagraphs) > 0 {
		sections = append(sections, fmt.Sprintf("## Part %03d\n%s", sectionNo, strings.Join(currentParagraphs, "\n")))
	}

	return strings.Join(sections, "\n\n")
}
