package service

import (
	"encoding/json"
	"net/http"
	"ping_pong_championship/commons"
)

const (
	errDecodingMsg = "Invalid request"
	errCode        = 404
)

type joinRequest struct {
	ServerURL string `json:"server_url"`
	Name      string `json:"name"`
}

func Join(w http.ResponseWriter, r *http.Request) {
	var req joinRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		commons.ResponseError(w, errDecodingMsg, errCode)
		return
	}

	id, err := joinGame(req.ServerURL, req.Name)
	if err != nil {
		commons.ResponseError(w, err.Error(), errCode)
		return
	}

	commons.ResponseJSON(w, map[string]int{"player_id": id})
}
