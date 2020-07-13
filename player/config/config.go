package config

import (
	"flag"
	"fmt"
	"log"
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
	port       string
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

func GetHostURL() string {
	return "localhost:" + port
}

func GetPort() string {
	return port
}

func ParseFlags() {
	portFlag := flag.String("p", "", "Port to start on")
	nameFlag := flag.String("n", "sample", "Name of the player")
	defenceLenFlag := flag.Int("d", 5, "Player's defence length")

	flag.Parse()

	name = *nameFlag
	defenceLen = *defenceLenFlag
	port = *portFlag

	if port == "" {
		log.Panic("port is empty")
	}
}

func GetHttpServer() *http.Server {
	if server != nil {
		return server
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: nil,
	}

	return server
}
