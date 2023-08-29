package util

import (
	"fmt"
	bcd "github.com/johnsonjh/gobcd"
	"io"
	"unicode/utf8"
)

var ErrReservedSpaceInsufficient = fmt.Errorf("insufficient reserved space")

// WriteText writes text to the writer using the Gen 1 US English character set.
// Once text is written, a terminator byte 0x50 is written.
// Additional padding of 0x00 bytes is written to ensure the entire reserved space is utilised.
// Should the reserved space be insufficient to write the text and terminator, an ErrReservedSpaceInsufficient error is returned.
func WriteText(w io.Writer, text string, reservedSpace int) error {
	usedSpace := utf8.RuneCount([]byte(text)) + 1 // +1 for terminator
	unusedSpace := reservedSpace - usedSpace
	if unusedSpace < 0 {
		return fmt.Errorf("cannot fit text %q in %d bytes: %w", text, reservedSpace, ErrReservedSpaceInsufficient)
	}

	const terminator, padding = 0x50, 0x00

	charConverter := map[rune]byte{
		'A': 0x80, 'B': 0x81, 'C': 0x82, 'D': 0x83, 'E': 0x84, 'F': 0x85,
		'G': 0x86, 'H': 0x87, 'I': 0x88, 'J': 0x89, 'K': 0x8A, 'L': 0x8B,
		'M': 0x8C, 'N': 0x8D, 'O': 0x8E, 'P': 0x8F, 'Q': 0x90, 'R': 0x91,
		'S': 0x92, 'T': 0x93, 'U': 0x94, 'V': 0x95, 'W': 0x96, 'X': 0x97,
		'Y': 0x98, 'Z': 0x99, '(': 0x9A, ')': 0x9B, ':': 0x9C, ';': 0x9D,
		'[': 0x9E, ']': 0x9F,

		'a': 0xA0, 'b': 0xA1, 'c': 0xA2, 'd': 0xA3, 'e': 0xA4, 'f': 0xA5,
		'g': 0xA6, 'h': 0xA7, 'i': 0xA8, 'j': 0xA9, 'k': 0xAA, 'l': 0xAB,
		'm': 0xAC, 'n': 0xAD, 'o': 0xAE, 'p': 0xAF, 'q': 0xB0, 'r': 0xB1,
		's': 0xB2, 't': 0xB3, 'u': 0xB4, 'v': 0xB5, 'w': 0xB6, 'x': 0xB7,
		'y': 0xB8, 'z': 0xB9,

		//'PK': 0x??, 'MN': 0x??,

		'-': 0xE3,
		'?': 0xE6, '!': 0xE7, '.': 0xE8,
		'♂': 0xEF,
		'/': 0xF3, ',': 0xF4, '♀': 0xF5,
	}

	for _, r := range text {
		char, present := charConverter[r]
		if !present {
			return fmt.Errorf("character %q is not available in the character set", string(r))
		}

		_, err := w.Write([]byte{char})
		if err != nil {
			return fmt.Errorf("failed to write rune %q as %b to writer: %w", r, char, err)
		}
	}

	_, err := w.Write([]byte{terminator})
	if err != nil {
		return fmt.Errorf("failed to write terminator: %w", err)
	}


	for i := 0; i < unusedSpace; i++ {
		_, err := w.Write([]byte{padding})
		if err != nil {
			return fmt.Errorf("failed to write padding: %w", err)
		}
	}

	return nil
}

func WriteBinaryCodedDecimal(w io.Writer, value uint64, reservedSpace int) error {
	b := bcd.FromUint(value, reservedSpace)
	_, err := w.Write(b)
	if err != nil {
		return fmt.Errorf("failed to write binary coded decimal: %w", err)
	}

	return err
}
