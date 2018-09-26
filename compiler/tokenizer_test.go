package compiler

import (
	"bytes"
	"log"
	"testing"
)

func TestCompiler(t *testing.T) {
	test := []byte("(print \"hello world\" 99)")

	r := bytes.NewReader(test)

	tokenizer := &Tokenizer{
		r,
	}

	for {
		token, err := tokenizer.ReadToken()
		if err != nil {
			log.Fatal(err)
		}

		log.Print(token)
	}
}
