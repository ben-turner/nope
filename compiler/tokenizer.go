package compiler

import (
	"errors"
	"fmt"
	"io"
)

type Tokenizer struct {
	r       io.ByteScanner
	column  int
	line    int
	current locatableByte
}

type TokenType uint8

const (
	Delimiter TokenType = iota
	Identifier
	StringLiteral
	IntegerLiteral
	FloatLiteral
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
	case IntegerLiteral:
		return "Integer"
	case FloatLiteral:
		return "Float"
	case BoolLiteral:
		return "Bool"
	default:
		return "Unknown"
	}
}

type Token struct {
	tokenType TokenType
	value     interface{}
	line      int
	col       int
}

func (t *Token) String() string {
	return fmt.Sprintf("%s %v", t.tokenType, t.value)
}

func readString(t *Tokenizer) (*Token, error) {
	str := []byte{}
	line := -1
	col := -1

	for {
		b, err := t.ReadByte()
		if err != nil {
			return nil, err
		}
		if line == -1 || col == -1 {
			line, col = t.line, t.col
		}

		if b == byte('"') {
			return &Token{
				tokenType: StringLiteral,
				value:     string(str),
				line:      line,
				col:       col,
			}, nil
		}

		if b == byte('\\') {
			n, err := t.ReadByte()
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
	b, err := t.ReadByte()
	if err != nil {
		return nil, err
	}
	line, col := t.line, t.col

	if !isLetter(b) {
		return nil, errors.New("Unexpected character")
	}

	ident := []byte{b}

	for {
		b, err = t.ReadByte()
		if err != err {
			return nil, err
		}

		if !isAlphaNum(b) {
			err := t.r.UnreadByte()
			return &Token{
				tokenType: Identifier,
				value:     string(ident),
				line:      line,
				col:       col,
			}, err
		}

		ident = append(ident, b)
	}
}

func readNumber(t *Tokenizer) (*Token, error) {
	value := 0
	line, col := -1, -1

	for {
		b, err := t.ReadByte()
		if err != nil {
			return nil, err
		}
		if line == -1 || col == -1 {
			line, col = t.line, t.col
		}

		switch true {
		case isNumber(b):
			value *= 10
			value += getNumberValue(b)
		case b == byte('.'):
			// decimal
		default:
			err := t.r.UnreadByte()
			if err != nil {
				return nil, err
			}
			return &Token{
				tokenType: IntegerLiteral,
				value:     value,
			}, nil
		}
	}
}

func (t *Tokenizer) ReadByte() (locatableByte, error) {
	b, err := t.r.ReadByte()
	if err != nil {
		return b, err
	}

	lb := locatableByte{
		value:  b,
		line:   t.line,
		column: t.column,
	}

	defer func() { t.current = lb }()

	if t.current == byte('\n') {
		t.column = 0
		t.line++
		return lb, nil
	}

	t.col++
	return b, err
}

func (t *Tokenizer) UnreadByte() error {
	err := t.r.UnreadByte()
	if err != nil {
		return err
	}

	t.col--
	return nil
}

func (t *Tokenizer) ReadToken() (*Token, error) {
	b, err := t.ReadByte()
	if err != nil {
		return nil, err
	}

	switch b {
	case byte('('), byte(')'):
		return &Token{
			tokenType: Delimiter,
			value:     string([]byte{b}),
			line:      t.line,
			col:       t.col,
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
