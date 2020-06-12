package referee

import (
	"encoding/json"
	"log"
	"net/http"

	. "ping_pong_championship/referee/models"
	. "ping_pong_championship/referee/services"

	"github.com/gorilla/mux"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	res, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}

// func joinHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hi there, I love %s!\n", r.URL.Path[1:])
// }

func joinHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var joinRequest JoinRequest
	err := json.NewDecoder(r.Body).Decode(&joinRequest)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var joinResponse JoinResponse
	joinResponse.PlayerId := services.join(joinRequest.Host)
	respondWithJson(w, http.StatusOK, joinResponse)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/join", joinHandler).Methods("POST")
	log.Println("Listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
