package service

import (
	"errors"
	"math/rand"
	"ping_pong_championship/player/config"
	"ping_pong_championship/player/domain"
	"time"
)

var (
	gamesMap map[string]*domain.Game
)

const maxRandomNo = 10

// AddGame ....
func AddGame(gameID string) {
	gamesMap[gameID] = &domain.Game{ID: gameID, DefenceArray: make([]int, config.GetDefenceLen())}
}

// DeleteGame ....
func DeleteGame(gameID string) {
	delete(gamesMap, gameID)
}

// GetDefenceNos ....
func GetDefenceNos(gameID string) ([]int, error) {
	game, exists := gamesMap[gameID]

	if !exists {
		return nil, errors.New("Invalid gameID " + gameID)
	}

	source := rand.NewSource(time.Now().UnixNano())
	randomGen := rand.New(source)

	defenceNos := game.DefenceArray

	for index := range defenceNos {
		// randomGen generates nos from 0 to maxRandomNo - 1
		// as defenceNos goes from 1 to maxRandomNo, so +1 at the end
		defenceNos[index] = randomGen.Intn(maxRandomNo) + 1
	}

	return defenceNos, nil
}

// GetRandomNumber ....
func GetRandomNumber(gameID string) (int, error) {
	game, exists := gamesMap[gameID]

	if !exists {
		return 0, errors.New("Invalid gameID " + gameID)
	}

	source := rand.NewSource(time.Now().UnixNano())
	randomGen := rand.New(source)

	randomNO := randomGen.Intn(len(game.DefenceArray))

	return game.DefenceArray[randomNO], nil
}
