package transport

import (
	"encoding/json"
	"log"
	"net/http"
	"ping_pong_championship/commons"
	"ping_pong_championship/commons/client"
	"ping_pong_championship/player/config"
	"ping_pong_championship/player/service"

	"github.com/gorilla/mux"
)

const (
	gameURLParam = "game_id"
	refereeURL   = "http://localhost:8080"
	errCode      = 404
)

func MakeHandler() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/game/{"+gameURLParam+"}", addGame).Methods(http.MethodPut)
	r.HandleFunc("/game/{"+gameURLParam+"}", deleteGame).Methods(http.MethodDelete)
	r.HandleFunc("/shutdown", shutDown).Methods(http.MethodDelete)
	r.HandleFunc("/game/{"+gameURLParam+"}/defence_numbers", getDefenceNos).Methods(http.MethodGet)
	r.HandleFunc("/game/{"+gameURLParam+"}/random_number", getRandomNoFromDefenceNos).Methods(http.MethodGet)
	return r
}

func addGame(w http.ResponseWriter, r *http.Request) {
	gameID := mux.Vars(r)[gameURLParam]
	service.AddGame(gameID)
}

func shutDown(w http.ResponseWriter, r *http.Request) {
	log.Fatal("ShutDown bye bye")
	// Hacky way to shutdown server
	// err := config.GetHttpServer().Shutdown(context.Background())
	// commons.HandleIfError("shutDown failed", err)
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

func JoinWithRefree() (err error) {

	type player struct {
		ServerURL string `json:"server_url"`
		Name      string `json:"name"`
	}

	playerInfo := player{ServerURL: config.GetHostURL(), Name: config.GetPlayerName()}

	requestData, err := json.Marshal(playerInfo)
	if err != nil {
		log.Printf("joinWithRefree - %s", err.Error())
		return
	}

	res, err := client.DoRequest(http.MethodPost, refereeURL+"/join", string(requestData))
	if err != nil {
		log.Printf("joinWithRefree - DoRequest %s %s ", refereeURL, err.Error())
		return
	}

	res.Body.Close()
	return
}
