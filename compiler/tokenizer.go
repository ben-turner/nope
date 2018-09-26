package compiler

import (
	"errors"
	"fmt"
	"io"
)

type Tokenizer struct {
	r io.ByteScanner
}

type TokenType uint8

const (
	Delimiter TokenType = iota
	Identifier
	StringLiteral
	NumberLiteral
	BoolLiteral
)

func (t TokenType) String() string {
	switch t {
	case Delimiter:
		return "Delimiter"
	case Identifier:
		return "Identifier"
	case StringLiteral:
		return "String"
	case NumberLiteral:
		return "Number"
	case BoolLiteral:
		return "Bool"
	default:
		return "Unknown"
	}
}

type Token struct {
	tokenType TokenType
	value     interface{}
}

func (t *Token) String() string {
	return fmt.Sprintf("%s %v", t.tokenType, t.value)
}

func isLetter(b byte) bool {
	return b == byte('a') ||
		b == byte('b') ||
		b == byte('c') ||
		b == byte('d') ||
		b == byte('e') ||
		b == byte('f') ||
		b == byte('g') ||
		b == byte('h') ||
		b == byte('i') ||
		b == byte('j') ||
		b == byte('k') ||
		b == byte('l') ||
		b == byte('m') ||
		b == byte('n') ||
		b == byte('o') ||
		b == byte('p') ||
		b == byte('q') ||
		b == byte('r') ||
		b == byte('s') ||
		b == byte('t') ||
		b == byte('u') ||
		b == byte('v') ||
		b == byte('w') ||
		b == byte('x') ||
		b == byte('y') ||
		b == byte('z')
}

func isNumber(b byte) bool {
	return b == byte('0') ||
		b == byte('1') ||
		b == byte('2') ||
		b == byte('3') ||
		b == byte('4') ||
		b == byte('5') ||
		b == byte('6') ||
		b == byte('7') ||
		b == byte('8') ||
		b == byte('9')
}

func getNumberValue(b byte) int {
	switch b {
	case byte('0'):
		return 0
	case byte('1'):
		return 1
	case byte('2'):
		return 2
	case byte('3'):
		return 3
	case byte('4'):
		return 4
	case byte('5'):
		return 5
	case byte('6'):
		return 6
	case byte('7'):
		return 7
	case byte('8'):
		return 8
	case byte('9'):
		return 9
	default:
		return -1
	}
}

func isAlphaNum(b byte) bool {
	return isNumber(b) || isLetter(b)
}

func readString(t *Tokenizer) (*Token, error) {
	str := []byte{}

	for {
		b, err := t.r.ReadByte()
		if err != nil {
			return nil, err
		}

		if b == byte('"') {
			return &Token{
				tokenType: StringLiteral,
				value:     string(str),
			}, nil
		}

		if b == byte('\\') {
			n, err := t.r.ReadByte()
			if err != nil {
				return nil, err
			}

			switch n {
			case byte('"'), byte('\\'):
				str = append(str, n)
			case byte('n'):
				str = append(str, byte('\n'))
			case byte('r'):
				str = append(str, byte('\r'))
			case byte('t'):
				str = append(str, byte('\t'))
			default:
				return nil, errors.New("Unexpected escape sequence")
			}

			continue
		}

		str = append(str, b)
	}
}

func readIdent(t *Tokenizer) (*Token, error) {
	b, err := t.r.ReadByte()
	if err != nil {
		return nil, err
	}

	if !isLetter(b) {
		return nil, errors.New("Unexpected character")
	}

	ident := []byte{b}

	for {
		b, err = t.r.ReadByte()
		if err != err {
			return nil, err
		}

		if !isAlphaNum(b) {
			err := t.r.UnreadByte()
			return &Token{
				tokenType: Identifier,
				value:     string(ident),
			}, err
		}

		ident = append(ident, b)
	}
}

func readNumber(t *Tokenizer) (*Token, error) {
	value := 0

	for {
		b, err := t.r.ReadByte()
		if err != nil {
			return nil, err
		}

		switch true {
		case isNumber(b):
			value *= 10
			value += getNumberValue(b)
		default:
			err := t.r.UnreadByte()
			if err != nil {
				return nil, err
			}
			return &Token{
				tokenType: NumberLiteral,
				value:     value,
			}, nil
		}
	}
}

func (t *Tokenizer) ReadToken() (*Token, error) {
	b, err := t.r.ReadByte()
	if err != nil {
		return nil, err
	}

	switch b {
	case byte('('), byte(')'):
		return &Token{
			tokenType: Delimiter,
			value:     string([]byte{b}),
		}, nil

	case byte(' '), byte('\t'), byte('\n'), byte('\r'):
		return t.ReadToken()

	case byte('"'):
		return readString(t)

	case byte('0'), byte('1'), byte('2'), byte('3'), byte('4'), byte('5'), byte('6'), byte('7'), byte('8'), byte('9'):
		err := t.r.UnreadByte()
		if err != nil {
			return nil, err
		}
		return readNumber(t)

	default:
		err := t.r.UnreadByte()
		if err != nil {
			return nil, err
		}
		return readIdent(t)
	}
}

// (func arg1 arg2)
