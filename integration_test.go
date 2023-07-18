// +build integration

package main_test

import (
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
)

var pokegen *dockertest.Resource

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	pokegen, err = pool.Run("pokegen", "latest", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("http://localhost:%s/health", pokegen.GetPort("8080/tcp")),
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to perform request: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("got unexpected status code: %d", resp.StatusCode)
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not get healthy instance: %s", err)
	}

	code := m.Run()

	if err := pool.Purge(pokegen); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestIntegrationHappy(t *testing.T) {
	assert := assert.New(t)
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://localhost:%s/gen", pokegen.GetPort("8080/tcp")),
		strings.NewReader(`{"player_name": "Red", "rival_name": "Gary"}`),
	)
	assert.NoError(err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(err)

	assert.Len(body, 32768)

	// pn = player name
	pnOffset, pnSize := 0x2598, 0xB
	assert.Equal(
		// "Red" + terminator + padding
		[]byte{0x91, 0xA4, 0xA3, 0x50, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		body[pnOffset:pnOffset+pnSize],
		"player name is incorrect",
	)

	// rn = rival name
	rnOffset, rnSize := 0x25F6, 0xB
	assert.Equal(
		// "Gary" + terminator + padding
		[]byte{0x86, 0xA0, 0xB1, 0xB8, 0x50, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		body[rnOffset:rnOffset+rnSize],
		"rival name is incorrect",
	)

	// cs = checksum
	csOffset, csSize := 0x3523, 0x1
	assert.Equal([]byte{0xC2}, body[csOffset:csOffset+csSize], "checksum is incorrect")

}
