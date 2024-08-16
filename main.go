package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"pokegen/internal/pokegen"
	"runtime"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	l = l.With(
		slog.String("app", "pokegen"),
		slog.String("go_version", runtime.Version()),
		slog.String("go_os", runtime.GOOS),
		slog.String("go_arch", runtime.GOARCH),
	)

	l.Info("Starting server")

	http.HandleFunc("POST /gen", genFile(l))
	http.HandleFunc("GET /health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			panic(fmt.Errorf("failed to write OK to /health"))
		}
	})
	return http.ListenAndServe(":8080", nil)
}

func genFile(l *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		l.Info("Received request",
			slog.String("method", req.Method),
			slog.String("path", req.URL.Path),
		)
		type schema struct {
			PlayerName string `json:"player_name"`
			RivalName  string `json:"rival_name"`
			Money      uint64 `json:"money"`
		}

		// Default values
		reqBody := schema{
			PlayerName: "RED",
			RivalName:  "BLUE",
			Money:      3000,
		}

		if req.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type is unsupported", http.StatusUnsupportedMediaType)
			return
		}

		err := json.NewDecoder(req.Body).Decode(&reqBody)
		if err != nil && err != io.EOF {
			var syntaxErr *json.SyntaxError
			if errors.As(err, &syntaxErr) {
				http.Error(w, fmt.Sprintf("syntax error at byte offset %d", syntaxErr.Offset), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := new(bytes.Buffer)
		_, err = pokegen.Gen(data, reqBody.PlayerName, reqBody.RivalName, reqBody.Money)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		if _, err = w.Write(data.Bytes()); err != nil {
			panic(err)
		}
	}
}
