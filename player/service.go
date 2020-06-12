package player

import "net/http"

type GameInfo struct {
	OpponentPlayerID int `json:"opponent_player_id"`
	PlayOrder        int `json:"play_order"`
}

type RandomNumberInfo struct {
	RandomNumber int `json:"random_number"`
}

type DefenceNumbersResponse struct {
	DefenceNumber []int `json:"defence_numbers"`
}

func (p *player) initGame(w http.ResponseWriter, r *http.Request) {

}
