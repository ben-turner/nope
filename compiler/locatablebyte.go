package tokenizer

import (
	"math"
)

type locatableByte struct {
	line   int
	column int
	value  byte
}

func (b *locatableByte) isNumber() bool {
	return b.value == byte('0') ||
		b.value == byte('1') ||
		b.value == byte('2') ||
		b.value == byte('3') ||
		b.value == byte('4') ||
		b.value == byte('5') ||
		b.value == byte('6') ||
		b.value == byte('7') ||
		b.value == byte('8') ||
		b.value == byte('9')
}

func (b *locatableByte) isLetter() bool {
	return b.value == byte('a') ||
		b.value == byte('b') ||
		b.value == byte('c') ||
		b.value == byte('d') ||
		b.value == byte('e') ||
		b.value == byte('f') ||
		b.value == byte('g') ||
		b.value == byte('h') ||
		b.value == byte('i') ||
		b.value == byte('j') ||
		b.value == byte('k') ||
		b.value == byte('l') ||
		b.value == byte('m') ||
		b.value == byte('n') ||
		b.value == byte('o') ||
		b.value == byte('p') ||
		b.value == byte('q') ||
		b.value == byte('r') ||
		b.value == byte('s') ||
		b.value == byte('t') ||
		b.value == byte('u') ||
		b.value == byte('v') ||
		b.value == byte('w') ||
		b.value == byte('x') ||
		b.value == byte('y') ||
		b.value == byte('z')
}

func (b *locatableByte) isAlphanumeric() bool {
	return b.isByte() || b.isLetter()
}

func (b *locatableByte) getNumericValue() int {
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
		return math.NaN
	}
}
