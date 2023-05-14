package pokegen_test

import (
	"github.com/stretchr/testify/assert"
	"os"
	"pokegen/internal/pokegen"
	"testing"
)

func TestGenNames(t *testing.T) {
	b, err := pokegen.Gen("JAM", "BOB", 3000)
	assert.NoError(t, err)

	dat, err := os.ReadFile("./jam-bob.sav")
	assert.NoError(t, err)

	assert.Equal(t, dat, b)
}

func TestGenNames2(t *testing.T) {
	b, err := pokegen.Gen("AaBbCcDdEe", "DEF", 3000)
	assert.NoError(t, err)

	dat, err := os.ReadFile("./abc-def.sav")
	assert.NoError(t, err)

	assert.Equal(t, dat, b)
}
