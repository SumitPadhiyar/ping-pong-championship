package commons

import (
	"encoding/json"
	"log"
	"net/http"
)

func ResponseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func ResponseError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	ResponseJSON(w, map[string]string{"error": message})
}

func HandleIfError(methodName string, err error) {
	if err == nil {
		return
	}

	log.Panicf("%s - %s\n", methodName, err.Error())
}
