package processors

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
)

type Sha512Processor struct{}

func (processor *Sha512Processor) Encode(ctx context.Context, url string) string {
	hasher := sha512.New()

	hasher.Write([]byte(url))
	hash := hasher.Sum(nil)

	encodedString := hex.EncodeToString(hash)

	return encodedString
}

func NewSha512Processor() *Sha512Processor {
	return &Sha512Processor{}
}
