package compiler

import (
	"bytes"
	"log"
	"testing"
)

func TestCompiler(t *testing.T) {
	program := `
	  (print
		  "hello world"
			99
			123.45
		)`
	test := []byte(program)

	r := bytes.NewReader(test)

	tokenizer := &Tokenizer{
		r,
		0,
		0,
		0,
	}

	for {
		token, err := tokenizer.ReadToken()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("%d:%d %v", tokenizer.line, tokenizer.col, token.String())
	}
}
