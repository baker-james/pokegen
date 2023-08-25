//go:build integration

package main_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func healthCheckCondition() bool {
	performHealthCheck := func() error {
		req, err := http.NewRequest(
			http.MethodGet,
			"http://localhost:8080/health",
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
	}

	if err := performHealthCheck(); err != nil {
		fmt.Println("health check failed:", err)
		return false
	}

	return true
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestIntegration_AcceptPlayerAndRivalNames(t *testing.T) {
	assert := assert.New(t)
	assert.Eventually(healthCheckCondition, 5*time.Second, 100*time.Millisecond)

	req, err := http.NewRequest(
		http.MethodGet,
		"http://localhost:8080/gen",
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

func TestIntegration_AcceptMoney(t *testing.T) {
	assert := assert.New(t)
	assert.Eventually(healthCheckCondition, 5*time.Second, 100*time.Millisecond)

	req, err := http.NewRequest(
		http.MethodGet,
		"http://localhost:8080/gen",
		strings.NewReader(`{"money": 4000}`),
	)
	assert.NoError(err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(err)

	assert.Len(body, 32768)

	defaultPlayerName := []byte{0x91, 0x84, 0x83}
	// pn = player name
	pnOffset, pnSize := 0x2598, 0xB
	assert.Equal(
		// "RED" + terminator + padding
		append(defaultPlayerName, 0x50, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00),
		body[pnOffset:pnOffset+pnSize],
		"player name is incorrect",
	)

	defaultRivalName := []byte{0x81, 0x8B, 0x94, 0x84}
	// rn = rival name
	rnOffset, rnSize := 0x25F6, 0xB
	assert.Equal(
		// "Gary" + terminator + padding
		append(defaultRivalName, 0x50, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00),
		body[rnOffset:rnOffset+rnSize],
		"rival name is incorrect",
	)

	moneyOffset, moneySize := 0x25F3, 0x3
	assert.Equal(
		[]byte{0x00, 0x40, 0x00},
		body[moneyOffset:moneyOffset+moneySize],
	)

	// cs = checksum
	csOffset, csSize := 0x3523, 0x1
	assert.Equal([]byte{0x5D}, body[csOffset:csOffset+csSize], "checksum is incorrect")
}

func TestIntegration_EmptyBody(t *testing.T) {
	assert := assert.New(t)
	assert.Eventually(healthCheckCondition, 5*time.Second, 100*time.Millisecond)

	req, err := http.NewRequest(
		http.MethodGet,
		"http://localhost:8080/gen",
		strings.NewReader(``),
	)
	assert.NoError(err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(err)

	assert.Len(body, 32768)

	defaultPlayerName := []byte{0x91, 0x84, 0x83}
	// pn = player name
	pnOffset, pnSize := 0x2598, 0xB
	assert.Equal(
		// "RED" + terminator + padding
		append(defaultPlayerName, 0x50, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00),
		body[pnOffset:pnOffset+pnSize],
		"player name is incorrect",
	)

	defaultRivalName := []byte{0x81, 0x8B, 0x94, 0x84}
	// rn = rival name
	rnOffset, rnSize := 0x25F6, 0xB
	assert.Equal(
		// "Gary" + terminator + padding
		append(defaultRivalName, 0x50, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00),
		body[rnOffset:rnOffset+rnSize],
		"rival name is incorrect",
	)

	// cs = checksum
	csOffset, csSize := 0x3523, 0x1
	assert.Equal([]byte{0x6D}, body[csOffset:csOffset+csSize], "checksum is incorrect")
}

func TestIntegration_InvalidBody(t *testing.T) {
	assert := assert.New(t)
	assert.Eventually(healthCheckCondition, 5*time.Second, 100*time.Millisecond)

	req, err := http.NewRequest(
		http.MethodGet,
		"http://localhost:8080/gen",
		strings.NewReader(`{this is invalid json}`),
	)
	assert.NoError(err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(err)
	assert.Equal("syntax error at byte offset 2\n", string(body))
}
