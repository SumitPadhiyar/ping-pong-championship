package player

import "fmt"

const (
	gamePath           = "/game"
	shutDown           = "/shutdown"
	gameDefenceNumbers = "/game/%s/defence_numbers"
	gameRandomNumbers  = "/game/%s/random_number"
)

func GetGamePUTPath(gameID string) string {
	return gamePath + "/" + gameID
}

func GetShutDownPath() string {
	return shutDown
}

func GetDeletePath(gameID string) string {
	return gamePath + "/" + gameID
}

func GetDefenceNumPath(gameID string) string {
	return fmt.Sprintf(gameDefenceNumbers, gameID)
}

func GetRandomPath(gameID string) string {
	return fmt.Sprintf(gameRandomNumbers, gameID)
}
