package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func GetTraderVersion(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(TraderVersionResponse{
		Version:  strings.Split(config.TraderImage,":")[1],
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}