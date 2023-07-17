package util_test

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"pokegen/internal/util"
	"testing"
)

type mockWriter func([]byte) (int, error)

func (m mockWriter) Write(p []byte) (n int, err error) {
	return m(p)
}

func TestWriteText_ExactSpaceForTextAndTerminator(t *testing.T) {
	buf := new(bytes.Buffer)
	err := util.WriteText(buf, "RED", 4)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x91, 0x84, 0x83, 0x50}, buf.Bytes())
}

func TestWriteText_AdditionalPaddingRequired(t *testing.T) {
	buf := new(bytes.Buffer)
	err := util.WriteText(buf, "RED", 6)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x91, 0x84, 0x83, 0x50, 0x00, 0x00}, buf.Bytes())
}

func TestWriteText_NoSpaceForText(t *testing.T) {
	buf := new(bytes.Buffer)
	err := util.WriteText(buf, "RED", 2)
	assert.ErrorIs(t, err, util.ErrReservedSpaceInsufficient)
}

func TestWriteText_NoSpaceForTerminator(t *testing.T) {
	buf := new(bytes.Buffer)
	err := util.WriteText(buf, "", 0)
	assert.ErrorIs(t, err, util.ErrReservedSpaceInsufficient)
}

func TestWriteText_ExactSpaceForTerminatorOnly(t *testing.T) {
	buf := new(bytes.Buffer)
	err := util.WriteText(buf, "", 1)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x50}, buf.Bytes())
}

func TestWriteText_WriteError(t *testing.T) {
	expectedErr := errors.New("expected")
	var mw mockWriter = func(_ []byte) (int, error) {
		return 0, expectedErr
	}

	err := util.WriteText(mw, "", 1)
	assert.ErrorIs(t, err, expectedErr)
}

func TestWriteBinaryCodedDecimal_EnoughSpace(t *testing.T) {
	buf := new(bytes.Buffer)
	err := util.WriteBinaryCodedDecimal(buf, 3000, 3)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x00, 0x30, 0x00}, buf.Bytes())
}

func TestWriteBinaryCodedDecimal_WriteError(t *testing.T) {
	expectedErr := errors.New("expected")
	var mw mockWriter = func(_ []byte) (int, error) {
		return 0, expectedErr
	}

	err := util.WriteBinaryCodedDecimal(mw, 3000, 3)
	assert.ErrorIs(t, err, expectedErr)
}

