package pokegen

import (
	"bytes"
	"fmt"
	bcd "github.com/johnsonjh/gobcd"
	"io"
	"unicode/utf8"
)

func Gen(playerName, rivalName string, money int) ([]byte, error) {
	var bank0 bytes.Buffer
	err := writeStart(&bank0)
	if err != nil {
		return nil, fmt.Errorf("start: %w", err)
	}

	var bank1 bytes.Buffer
	err = writeMiddle(&bank1, playerName, rivalName, money)
	if err != nil {
		return nil, fmt.Errorf("middle: %w", err)
	}

	var bankn bytes.Buffer
	err = writeEnd(&bankn)
	if err != nil {
		return nil, fmt.Errorf("end: %w", err)
	}

	all := append(bank0.Bytes(), append(bank1.Bytes(), bankn.Bytes()...)...)

	return all, nil
}

func writeStart(w io.Writer) error {
	for i := 0; i < 132; i++ {
		_, err := w.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	_, err := w.Write([]byte{
		0x03, 0x0C, 0x10, 0x10, 0x20, 0x20,
		0x20, 0x10, 0x10, 0x08, 0x18, 0x20,
		0x20, 0x40, 0x48, 0x38, 0x08, 0x10,
		0x11, 0x11, 0x22, 0x22, 0x1C,
	})
	if err != nil {
		return fmt.Errorf("failed to write null byte: %w", err)
	}

	for i := 0; i < 33; i++ {
		_, err := w.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	_, err = w.Write([]byte{
		0xC0, 0x30, 0x08, 0x08, 0x04,
		0x04, 0x04, 0x08, 0x08, 0x10,
		0x18, 0x04, 0x04, 0x02, 0x12,
		0x1C, 0x10, 0x08, 0x88, 0x88,
		0x44, 0x44, 0x38,
	})
	if err != nil {
		return fmt.Errorf("failed to write null byte: %w", err)
	}

	for i := 0; i < 445; i++ {
		_, err := w.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	_, err = w.Write([]byte{
		0x03, 0x03, 0x0C, 0x0C, 0x10, 0x10,
		0x10, 0x10, 0x20, 0x20, 0x20, 0x20,
		0x20, 0x20, 0x10, 0x10, 0x10, 0x10,
		0x08, 0x08, 0x18, 0x18, 0x20, 0x20,
		0x20, 0x20, 0x40, 0x40, 0x48, 0x48,
		0x38, 0x38, 0x08, 0x08, 0x10, 0x10,
		0x11, 0x11, 0x11, 0x11, 0x22, 0x22,
		0x22, 0x22, 0x1C, 0x1C,
	})
	if err != nil {
		return fmt.Errorf("failed to write null byte: %w", err)
	}

	for i := 0; i < 66; i++ {
		_, err := w.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	_, err = w.Write([]byte{
		0xC0, 0xC0, 0x30, 0x30, 0x08, 0x08,
		0x08, 0x08, 0x04, 0x04, 0x04, 0x04,
		0x04, 0x04, 0x08, 0x08, 0x08, 0x08,
		0x10, 0x10, 0x18, 0x18, 0x04, 0x04,
		0x04, 0x04, 0x02, 0x02, 0x12, 0x12,
		0x1C, 0x1C, 0x10, 0x10, 0x08, 0x08,
		0x88, 0x88, 0x88, 0x88, 0x44, 0x44,
		0x44, 0x44, 0x38, 0x38,
	})
	if err != nil {
		return fmt.Errorf("failed to write null byte: %w", err)
	}

	for i := 0; i < 362; i++ {
		_, err := w.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 8448; i++ {
		_, err := w.Write([]byte{0xFF})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	return nil
}

func writeBinaryCodedDecimal(b *checksumWriter, i int) error {
	m := bcd.FromUint(uint64(i), 3)
	_, err := b.Write(m)
	if err != nil {
		return fmt.Errorf("failed to write null byte: %w", err)
	}

	return err
}

func writeUserInput(w io.Writer, text string, reservedSpace int) error {
	const terminator = 0x50

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

		w.Write([]byte{char})
	}

	w.Write([]byte{terminator})

	unusedSpace := reservedSpace - utf8.RuneCount([]byte(text))

	for i := 0; i < unusedSpace; i++ {
		w.Write([]byte{0x00})
	}

	return nil
}

type checksumWriter struct {
	w   io.Writer
	sum byte
}

func (csw *checksumWriter) Write(bytes []byte) (int, error) {
	for _, b := range bytes {
		csw.sum += b
	}
	return csw.w.Write(bytes)
}

func (csw checksumWriter) WriteChecksum() (int, error) {
	return csw.w.Write([]byte{^csw.sum})
}

func writeMiddle(w io.Writer, playerName, rivalName string, money int) error {
	var err error

	var csw = &checksumWriter{
		w: w,
	}

	const playerNameSpace = 10
	err = writeUserInput(csw, playerName, playerNameSpace)
	if err != nil {
		return fmt.Errorf("failed to write null byte: %w", err)
	}

	for i := 0; i < 39; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1; i++ {
		_, err := csw.Write([]byte{0xFF})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 40; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	err = writeBinaryCodedDecimal(csw, money)
	if err != nil {
		return fmt.Errorf("failed to write null byte: %w", err)
	}

	const rivalNameSpace = 10
	err = writeUserInput(csw, rivalName, rivalNameSpace)
	if err != nil {
		return fmt.Errorf("failed to write null byte: %w", err)
	}

	_, err = csw.Write([]byte{
		0x03, 0x00, 0x00, 0x01,
	})
	if err != nil {
		return fmt.Errorf("failed to write null byte: %w", err)
	}

	var playerID = []byte{0xC0, 0xB2}
	_, err = csw.Write(playerID)
	if err != nil {
		return fmt.Errorf("failed to write null byte: %w", err)
	}

	_, err = csw.Write([]byte{

		0xBA, 0x02,
		0x00, 0x26, 0x12, 0xC7, 0x06, 0x03, 0x00,
		0x01, 0x00, 0x00, 0x04, 0x04, 0x04, 0x10,
		0x40, 0xCF, 0x40, 0xB0, 0x40, 0x00, 0xFF,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0xFF, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0xFF, 0x00,
	})
	if err != nil {
		return fmt.Errorf("failed to write null byte: %w", err)
	}

	for i := 0; i < 20; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	_, err = csw.Write([]byte{
		0x00, 0xD0, 0x40, 0x00, 0x00,
		0x0A, 0x01, 0x01, 0x07, 0x02,
		0x25, 0x00,
	})
	if err != nil {
		return fmt.Errorf("failed to write null byte: %w", err)
	}

	for i := 0; i < 123; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1; i++ {
		_, err := csw.Write([]byte{0xFF})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 244; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	_, err = csw.Write([]byte{
		0x08, 0x08, 0x00, 0x98, 0x00, 0x08, 0x00,
		0x19, 0x70, 0x52, 0xE0, 0x4D, 0x49, 0x17,
		0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x14, 0x01, 0xFF, 0x00, 0x00,
	})
	if err != nil {
		return fmt.Errorf("failed to write null byte: %w", err)
	}

	for i := 0; i < 102; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	_, err = csw.Write([]byte{
		0xA5, 0x00, 0x7E, 0x01, 0x0C, 0x41, 0x02, 0x00, 0x10,
		0x10, 0x00, 0x00, 0x0C, 0x00, 0x02, 0x00, 0x80, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40,
		0x9E, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00,
	})
	if err != nil {
		return fmt.Errorf("failed to write null byte: %w", err)
	}

	for i := 0; i < 60; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1; i++ {
		_, err := csw.Write([]byte{0x01})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 270; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1; i++ {
		_, err := csw.Write([]byte{0xFF})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 22; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1; i++ {
		_, err := csw.Write([]byte{0x01})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 785; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	const seconds = 0x07
	const frame = 0x05

	_, err = csw.Write([]byte{
		seconds, frame,
	})
	if err != nil {
		return fmt.Errorf("failed to write null byte: %w", err)
	}

	for i := 0; i < 58; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	_, err = csw.Write([]byte{
		0x01, 0x00, 0xFF, 0x00, 0x3C, 0x00, 0x40, 0x00, 0x00, 0x04, 0x40, 0x40,
	})
	if err != nil {
		return fmt.Errorf("failed to write null byte: %w", err)
	}

	for i := 0; i < 6; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1; i++ {
		_, err := csw.Write([]byte{0xFF})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 15; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1; i++ {
		_, err := csw.Write([]byte{0xFF})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 15; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1; i++ {
		_, err := csw.Write([]byte{0xFF})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 15; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1; i++ {
		_, err := csw.Write([]byte{0xFF})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 15; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1; i++ {
		_, err := csw.Write([]byte{0xFF})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 15; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1; i++ {
		_, err := csw.Write([]byte{0xFF})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 15; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1; i++ {
		_, err := csw.Write([]byte{0xFF})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 15; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1; i++ {
		_, err := csw.Write([]byte{0xFF})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 15; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1; i++ {
		_, err := csw.Write([]byte{0xFF})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 15; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1; i++ {
		_, err := csw.Write([]byte{0xFF})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 15; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1; i++ {
		_, err := csw.Write([]byte{0xFF})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 15; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1; i++ {
		_, err := csw.Write([]byte{0xFF})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 15; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1; i++ {
		_, err := csw.Write([]byte{0xFF})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 15; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1; i++ {
		_, err := csw.Write([]byte{0xFF})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 15; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1; i++ {
		_, err := csw.Write([]byte{0xFF})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 26; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	_, err = csw.Write([]byte{
		0x01, 0x01,
	})
	if err != nil {
		return fmt.Errorf("failed to write null byte: %w", err)
	}

	for i := 0; i < 242; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1; i++ {
		_, err := csw.Write([]byte{0xFF})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 403; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1; i++ {
		_, err := csw.Write([]byte{0xFF})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	for i := 0; i < 1121; i++ {
		_, err := csw.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	_, err = csw.WriteChecksum()
	if err != nil {
		return fmt.Errorf("failed to write checksum: %w", err)
	}

	return nil
}

func writeEnd(w io.Writer) error {
	for i := 0; i < 19164; i++ {
		_, err := w.Write([]byte{0xFF})
		if err != nil {
			return fmt.Errorf("failed to write null byte: %w", err)
		}
	}

	return nil
}
