package processors

import (
	"context"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

type Base62Processor struct{}

func (processor *Base62Processor) Encode(ctx context.Context, series uint64) string {
	encoded := make([]byte, 7)

	for i := 6; i >= 0; i-- {
		index := series % 62
		encoded[i] = base62Chars[index]
		series /= 62
	}

	return string(encoded)
}

func NewBase62Processor() *Base62Processor {
	return &Base62Processor{}
}
