package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pokegen/internal/pokegen"
)

func main() {
	http.HandleFunc("/gen", genFile)
	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			panic(fmt.Errorf("failed to write OK to /health"))
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func genFile(w http.ResponseWriter, req *http.Request) {
	type reqSchema struct {
		PlayerName string `json:"player_name"`
		RivalName  string `json:"rival_name"`
	}
	var r reqSchema

	err := json.NewDecoder(req.Body).Decode(&r)
	if err != nil {
		panic(fmt.Errorf("unable to decode: %w", err))
	}

	data := new(bytes.Buffer)
	_, err = pokegen.Gen(data, r.PlayerName, r.RivalName, 3000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err = w.Write(data.Bytes()); err != nil {
		panic(err)
	}
}
