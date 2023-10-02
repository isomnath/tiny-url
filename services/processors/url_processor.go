package processors

import (
	"context"
	"net/url"
	"strings"
)

type URLProcessor struct{}

func (processor *URLProcessor) ExtractDomain(ctx context.Context, rawURL string) string {
	parsedURL, _ := url.Parse(rawURL)
	parts := strings.Split(parsedURL.Hostname(), ".")
	if len(parts) > 1 {
		return parts[len(parts)-2]
	}
	return parsedURL.Hostname()
}

func NewURLProcessor() *URLProcessor {
	return &URLProcessor{}
}
