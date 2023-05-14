package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pokegen/internal/pokegen"
)

func main() {
	http.HandleFunc("/gen", genFile)
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

	b, err := pokegen.Gen(r.PlayerName, r.RivalName, 3000)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(b)
	if err != nil {
		panic(err)
	}

	return
}
