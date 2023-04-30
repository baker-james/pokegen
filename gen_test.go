package pokegen_test

import (
	"github.com/stretchr/testify/assert"
	"os"
	"pokegen"
	"testing"
)

func TestGenNames(t *testing.T) {
	b, err := pokegen.Gen("JAM", "BOB")
	assert.NoError(t, err)

	dat, err := os.ReadFile("./jam-bob.sav")
	assert.NoError(t, err)

	assert.Equal(t, dat, b)

	err = os.WriteFile("./Pokemon Red.sav", b, 0644)
	if err != nil {
		panic(err)
	}
}


func TestGenNames2(t *testing.T) {
	b, err := pokegen.Gen("AaBbCcDdEe", "DEF")
	assert.NoError(t, err)

	dat, err := os.ReadFile("./abc-def.sav")
	assert.NoError(t, err)

	assert.Equal(t, dat, b)

	err = os.WriteFile("./Pokemon Red.sav", b, 0644)
	if err != nil {
		panic(err)
	}
}

