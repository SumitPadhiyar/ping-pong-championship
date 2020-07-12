package config

import (
	"flag"
	"fmt"
	"net/http"
)

const (
	gamePath           = "/game"
	shutDown           = "/shutdown"
	gameDefenceNumbers = "/game/%s/defence_numbers"
	gameRandomNumbers  = "/game/%s/random_number"
)

var (
	server     *http.Server
	defenceLen int
	name       string
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

func GetDefenceLen() int {
	return defenceLen
}

func GetPlayerName() string {
	return name
}

func ParseFlags() {
	nameFlag := flag.String("name", "sample", "Name of the player")
	defenceLenFlag := flag.Int("defence_length", 5, "Player's defence length")

	flag.Parse()

	name = *nameFlag
	defenceLen = *defenceLenFlag

}

func GetHttpServer() *http.Server {
	if server != nil {
		return server
	}

	server := &http.Server{
		Addr:    ":8081",
		Handler: nil,
	}

	return server
}
