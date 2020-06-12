package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"ping_pong_championship/client"
	"ping_pong_championship/player"
	"ping_pong_championship/referee/commons"
	"ping_pong_championship/referee/models"
	"strconv"
	"sync"
)

var (
	players        []models.Player
	games          []models.Game
	currentPlayers []*models.Player
	currentGames   []*models.Game
	mux            sync.Mutex
)

const (
	maxPlayer = 8
)

type playerInfoToSend struct {
	OpponentPlayerID int  `json:"opponent_player_id"`
	PlayFirst        bool `json:"play_order"`
}

type playerDefenceResponse struct {
	DefenceNumbers []int `json:"defence_numbers"`
}

func joinGame(url string, name string) (int, error) {

	mux.Lock()
	defer mux.Unlock()

	if len(players) >= maxPlayer {
		return 0, errors.New("Max player reached")
	}

	player := models.Player{ID: len(players), Name: name, PlayerRemoteURL: url}
	players = append(players, player)
	currentPlayers = append(currentPlayers, &player)

	if len(players) == maxPlayer {
		go initGame()
	}

	return player.ID, nil
}

func initGame() {

	currentGames = nil

	for i := 0; i <= (len(currentPlayers) / 2); i = i + 2 {

		player1 := models.GamePlayerInfo{}
		player1.Player = currentPlayers[i]

		player2 := models.GamePlayerInfo{}
		player2.Player = currentPlayers[i+1]

		game := models.Game{}
		game.ID = len(games)
		game.FirstPlayer = &player1
		game.SecondPlayer = &player2

		games = append(games, game)
		currentGames = append(currentGames, &game)

		go startGame(game)

	}

}

func startGame(game models.Game) {

	player1URL := "http;//" + game.FirstPlayer.Player.PlayerRemoteURL
	player2URL := "http;//" + game.SecondPlayer.Player.PlayerRemoteURL

	// Send game info to player 1

	playerInfo := playerInfoToSend{}
	playerInfo.OpponentPlayerID = game.SecondPlayer.Player.ID
	playerInfo.PlayFirst = true

	requestData, err := json.Marshal(playerInfo)
	commons.HandleIfError("startGame", err)

	res, err := client.DoRequest(http.MethodPut, player1URL+player.GetGamePUTPath(strconv.Itoa(game.ID)), string(requestData))
	res.Body.Close()
	commons.HandleIfError("startGame - GetGamePUTPath", err)

	// Send game info to player 2

	playerInfo.OpponentPlayerID = game.FirstPlayer.Player.ID
	playerInfo.PlayFirst = false

	requestData, err = json.Marshal(playerInfo)
	commons.HandleIfError("startGame", err)

	res, err = client.DoRequest(http.MethodPut, player2URL+player.GetGamePUTPath(strconv.Itoa(game.ID)), string(requestData))
	res.Body.Close()
	commons.HandleIfError("startGame - GetGamePUTPath", err)

	// Get defence numbers from player 1

	res, err = client.DoRequest(http.MethodGet, player1URL+player.GetDefenceNumPath(strconv.Itoa(game.ID)), "")
	defer res.Body.Close()
	commons.HandleIfError("startGame - GetDefenceNumPath", err)

	var defenceNos playerDefenceResponse
	err = json.NewDecoder(res.Body).Decode(&defenceNos)
	commons.HandleIfError("startGame - GetDefenceNumPath - Json decoding", err)
	game.FirstPlayer.DefenceArray = defenceNos.DefenceNumbers

	// Get defence numbers from player 2

	res, err = client.DoRequest(http.MethodGet, player2URL+player.GetDefenceNumPath(strconv.Itoa(game.ID)), "")
	defer res.Body.Close()
	commons.HandleIfError("startGame - GetDefenceNumPath", err)

	err = json.NewDecoder(res.Body).Decode(&defenceNos)
	commons.HandleIfError("startGame - GetDefenceNumPath - Json decoding", err)
	game.SecondPlayer.DefenceArray = defenceNos.DefenceNumbers

}
