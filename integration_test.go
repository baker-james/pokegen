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

const (
	playerNameStart = 0x2598
	playerNameEnd   = 0x25A3

	rivalNameStart = 0x25F6
	rivalNameEnd   = 0x2601

	checksumStart = 0x3523
	checksumEnd   = 0x3524
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
		http.MethodPost,
		"http://localhost:8080/gen",
		strings.NewReader(`{"player_name": "Red", "rival_name": "Gary"}`),
	)
	assert.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(err)

	assert.Len(body, 32768)

	assert.Equal(
		// "Red" + terminator + padding
		[]byte{0x91, 0xA4, 0xA3, 0x50, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		body[playerNameStart:playerNameEnd],
		"player name is incorrect",
	)

	assert.Equal(
		// "Gary" + terminator + padding
		[]byte{0x86, 0xA0, 0xB1, 0xB8, 0x50, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		body[rivalNameStart:rivalNameEnd],
		"rival name is incorrect",
	)

	assert.Equal([]byte{0xC2}, body[checksumStart:checksumEnd], "checksum is incorrect")
}

func TestIntegration_AcceptMoney(t *testing.T) {
	assert := assert.New(t)
	assert.Eventually(healthCheckCondition, 5*time.Second, 100*time.Millisecond)

	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/gen",
		strings.NewReader(`{"money": 4000}`),
	)
	assert.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(err)

	assert.Len(body, 32768)

	defaultPlayerName := []byte{0x91, 0x84, 0x83}
	assert.Equal(
		// "RED" + terminator + padding
		append(defaultPlayerName, 0x50, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00),
		body[playerNameStart:playerNameEnd],
		"player name is incorrect",
	)

	defaultRivalName := []byte{0x81, 0x8B, 0x94, 0x84}
	assert.Equal(
		// "Gary" + terminator + padding
		append(defaultRivalName, 0x50, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00),
		body[rivalNameStart:rivalNameEnd],
		"rival name is incorrect",
	)

	moneyOffset, moneySize := 0x25F3, 0x3
	assert.Equal(
		[]byte{0x00, 0x40, 0x00},
		body[moneyOffset:moneyOffset+moneySize],
	)

	assert.Equal([]byte{0x5D}, body[checksumStart:checksumEnd], "checksum is incorrect")
}

func TestIntegration_EmptyBody(t *testing.T) {
	assert := assert.New(t)
	assert.Eventually(healthCheckCondition, 5*time.Second, 100*time.Millisecond)

	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/gen",
		strings.NewReader(``),
	)
	assert.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(err)

	assert.Len(body, 32768)

	defaultPlayerName := []byte{0x91, 0x84, 0x83}
	assert.Equal(
		// "RED" + terminator + padding
		append(defaultPlayerName, 0x50, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00),
		body[playerNameStart:playerNameEnd],
		"player name is incorrect",
	)

	defaultRivalName := []byte{0x81, 0x8B, 0x94, 0x84}
	assert.Equal(
		// "Gary" + terminator + padding
		append(defaultRivalName, 0x50, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00),
		body[rivalNameStart:rivalNameEnd],
		"rival name is incorrect",
	)

	assert.Equal([]byte{0x6D}, body[checksumStart:checksumEnd], "checksum is incorrect")
}

func TestIntegration_InvalidBody(t *testing.T) {
	assert := assert.New(t)
	assert.Eventually(healthCheckCondition, 5*time.Second, 100*time.Millisecond)

	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/gen",
		strings.NewReader(`{this is invalid json}`),
	)
	assert.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusBadRequest, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(err)
	assert.Equal("syntax error at byte offset 2\n", string(body))
}

func TestIntegration_InvalidMethods(t *testing.T) {
	assert := assert.New(t)
	assert.Eventually(healthCheckCondition, 5*time.Second, 100*time.Millisecond)

	invalidMethods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}

	for _, method := range invalidMethods {
		req, err := http.NewRequest(
			method,
			"http://localhost:8080/gen",
			strings.NewReader(`{"money": 4000}`),
		)
		assert.NoError(err)

		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(err)
		assert.Equal(http.StatusMethodNotAllowed, resp.StatusCode)
	}
}

func TestIntegration_UnsupportedContentType(t *testing.T) {
	assert := assert.New(t)
	assert.Eventually(healthCheckCondition, 5*time.Second, 100*time.Millisecond)

	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/gen",
		strings.NewReader(`{}`),
	)
	assert.NoError(err)

	req.Header.Set("Content-Type", "application/xml")
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusUnsupportedMediaType, resp.StatusCode)
}

func TestIntegration_MissingContentType(t *testing.T) {
	assert := assert.New(t)
	assert.Eventually(healthCheckCondition, 5*time.Second, 100*time.Millisecond)

	req, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/gen",
		strings.NewReader(`{}`),
	)
	assert.NoError(err)

	//req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusUnsupportedMediaType, resp.StatusCode)
}
