package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"pokegen/internal/pokegen"
)

func main() {
	http.HandleFunc("POST /gen", genFile)
	http.HandleFunc("GET /health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			panic(fmt.Errorf("failed to write OK to /health"))
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func genFile(w http.ResponseWriter, req *http.Request) {
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
