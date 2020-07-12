package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"ping_pong_championship/commons"
	"ping_pong_championship/commons/client"
	"ping_pong_championship/player/config"
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
	winScore  = 5
)

type playerInfoToSend struct {
	OpponentPlayerID int  `json:"opponent_player_id"`
	PlayFirst        bool `json:"play_order"`
}

type playerDefenceResponse struct {
	DefenceNumbers []int `json:"defence_numbers"`
}

type playerAttackNumberResponse struct {
	AttackNumber int `json:"attack_number"`
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
		go initGames()
	}

	return player.ID, nil
}

func initGames() {

	var wg sync.WaitGroup

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

		// wait for next goroutine
		wg.Add(1)

		go startGame(game, &wg)

	}

	// Waiting for all goroutines
	wg.Wait()

	currentPlayers = nil

	// calculate currentPlayers
	for _, g := range currentGames {
		currentPlayers = append(currentPlayers, g.WinnerPlayer.Player)
	}

	if len(currentPlayers) == 1 {
		declareChampions()
	} else {
		initGames()
	}
}

func declareChampions() {

	if len(currentPlayers) != 1 {
		panic(errors.New("Invalid program state - len(currentPlayers) != 1 "))
	}

	fmt.Printf("The Champion is " + currentPlayers[0].Name)
	fmt.Printf("GameID\tFinal Score\tWinner\n")

	for _, g := range games {
		fmt.Printf("%d\t%d\t%s\n", g.ID, g.WinnerPlayer.Score, g.WinnerPlayer.Player.Name)
	}

}

func generateDefenceMap(defenceNos []int) map[int]struct{} {
	defenceMap := make(map[int]struct{}, len(defenceNos))
	var member struct{}

	for _, n := range defenceNos {
		defenceMap[n] = member
	}

	return defenceMap
}

func startGame(game models.Game, wg *sync.WaitGroup) {

	defer wg.Done()

	player1URL := "http://" + game.FirstPlayer.Player.PlayerRemoteURL
	player2URL := "http://" + game.SecondPlayer.Player.PlayerRemoteURL

	// Send game info to player 1

	playerInfo := playerInfoToSend{}
	playerInfo.OpponentPlayerID = game.SecondPlayer.Player.ID
	playerInfo.PlayFirst = true

	requestData, err := json.Marshal(playerInfo)
	commons.HandleIfError("startGame", err)

	res, err := client.DoRequest(http.MethodPut, player1URL+config.GetGamePUTPath(strconv.Itoa(game.ID)), string(requestData))
	res.Body.Close()
	commons.HandleIfError("startGame - GetGamePUTPath", err)

	// Send game info to player 2

	playerInfo.OpponentPlayerID = game.FirstPlayer.Player.ID
	playerInfo.PlayFirst = false

	requestData, err = json.Marshal(playerInfo)
	commons.HandleIfError("startGame", err)

	res, err = client.DoRequest(http.MethodPut, player2URL+config.GetGamePUTPath(strconv.Itoa(game.ID)), string(requestData))
	res.Body.Close()
	commons.HandleIfError("startGame - GetGamePUTPath", err)

	// Get defence numbers from player 1

	res, err = client.DoRequest(http.MethodGet, player1URL+config.GetDefenceNumPath(strconv.Itoa(game.ID)), "")
	commons.HandleIfError("startGame - GetDefenceNumPath", err)
	defer res.Body.Close()

	var defenceNos playerDefenceResponse
	err = json.NewDecoder(res.Body).Decode(&defenceNos)
	commons.HandleIfError("startGame - GetDefenceNumPath - Json decoding", err)
	game.FirstPlayer.DefenceMap = generateDefenceMap(defenceNos.DefenceNumbers)

	// Get defence numbers from player 2

	res, err = client.DoRequest(http.MethodGet, player2URL+config.GetDefenceNumPath(strconv.Itoa(game.ID)), "")
	commons.HandleIfError("startGame - GetDefenceNumPath", err)
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&defenceNos)
	commons.HandleIfError("startGame - GetDefenceNumPath - Json decoding", err)
	game.SecondPlayer.DefenceMap = generateDefenceMap(defenceNos.DefenceNumbers)

	// start game here
	playGame(game)

}

func playGame(game models.Game) {
	firstPlayer := game.FirstPlayer
	secondPlayer := game.SecondPlayer

	gameIDStr := strconv.Itoa(game.ID)

	player1URL := "http://" + firstPlayer.Player.PlayerRemoteURL + config.GetRandomPath(gameIDStr)
	player2URL := "http://" + secondPlayer.Player.PlayerRemoteURL + config.GetRandomPath(gameIDStr)

	offensivePlayer := firstPlayer
	offensivePlayerURL := player1URL
	defensivePlayer := secondPlayer
	defensivePlayerURL := player2URL

	for firstPlayer.Score < winScore && secondPlayer.Score < winScore {

		// Ask for randoom no. from offensive player
		res, err := client.DoRequest(http.MethodGet, offensivePlayerURL, "")
		commons.HandleIfError("playGame - GetRandomPath "+offensivePlayerURL, err)
		defer res.Body.Close()

		// Get attack no
		var attackNumber playerAttackNumberResponse
		err = json.NewDecoder(res.Body).Decode(&attackNumber)
		commons.HandleIfError("playGame - GetRandomPath "+offensivePlayerURL+" Json decoding", err)

		_, isNumberFound := defensivePlayer.DefenceMap[attackNumber.AttackNumber]

		if !isNumberFound {
			offensivePlayer.Score = offensivePlayer.Score + 1
		} else {
			defensivePlayer.Score = defensivePlayer.Score + 1

			// swap players
			temp := offensivePlayer
			offensivePlayer = defensivePlayer
			defensivePlayer = temp

			// swap players URL
			tempURL := offensivePlayerURL
			offensivePlayerURL = defensivePlayerURL
			defensivePlayerURL = tempURL
		}
	}

	var shutDownURL string

	if firstPlayer.Score == winScore {
		game.WinnerPlayer = firstPlayer
		shutDownURL = "http://" + firstPlayer.Player.PlayerRemoteURL + config.GetShutDownPath()
	} else {
		game.WinnerPlayer = secondPlayer
		shutDownURL = "http://" + secondPlayer.Player.PlayerRemoteURL + config.GetShutDownPath()
	}

	// Ask shutDown looser
	res, err := client.DoRequest(http.MethodPut, shutDownURL, "")
	fmt.Printf("%s - %s\n", "playGame - ShutDown "+shutDownURL, err.Error())
	defer res.Body.Close()

	deleteURL := "http://" + game.WinnerPlayer.Player.PlayerRemoteURL + config.GetDeletePath(strconv.Itoa(game.ID))

	// Delete game
	res, err = client.DoRequest(http.MethodDelete, deleteURL, "")
	commons.HandleIfError("playGame - Delete "+deleteURL, err)
	res.Body.Close()

}
