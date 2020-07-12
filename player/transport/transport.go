package transport

import (
	"context"
	"fmt"
	"net/http"
	"ping_pong_championship/commons"
	"ping_pong_championship/player/config"
	"ping_pong_championship/player/service"

	"github.com/gorilla/mux"
)

const (
	gameURLParam = "game_id"
	errCode      = 404
)

func MakeHandler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/game/{"+gameURLParam+"}", addGame).Methods(http.MethodPut)
	r.HandleFunc("/game/{"+gameURLParam+"}", deleteGame).Methods(http.MethodDelete)
	r.HandleFunc("/shutdown", shutDown).Methods(http.MethodPut)
	r.HandleFunc("/game/{"+gameURLParam+"}", getDefenceNos).Methods(http.MethodGet)
	r.HandleFunc("/game/{"+gameURLParam+"}", getRandomNoFromDefenceNos).Methods(http.MethodGet)
	return r
}

func addGame(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)[gameURLParam]
	service.AddGame(gameID)
}

func shutDown(w http.ResponseWriter, r *http.Request) {
	// Hacky way to shutdown server
	err := config.GetHttpServer().Shutdown(context.Background())
	if err != nil {
		fmt.Printf("shutDown failed with %s", err.Error())
	}
}

func deleteGame(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)[gameURLParam]
	service.DeleteGame(gameID)
}

func getDefenceNos(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)[gameURLParam]
	nos, err := service.GetDefenceNos(gameID)

	if err != nil {
		commons.ResponseError(w, err.Error(), errCode)
		return
	}

	commons.ResponseJSON(w, map[string][]int{"defence_numbers": nos})
}

func getRandomNoFromDefenceNos(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)[gameURLParam]
	randomNo, err := service.GetRandomNumber(gameID)

	if err != nil {
		commons.ResponseError(w, err.Error(), errCode)
		return
	}

	commons.ResponseJSON(w, map[string]int{"random_number": randomNo})
}
