package internal

import (
	"github.com/tiktoken-go/tokenizer"
)

func CountTokens(text string, encoding string) (int, error) {
	enc, err := tokenizer.Get(tokenizer.Encoding(encoding))
	if err != nil {
		return 0, err
	}

	tokens, _, err := enc.Encode(text)
	if err != nil {
		return 0, err
	}

	return len(tokens), nil
}