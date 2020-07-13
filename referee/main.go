package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/context"

	services "ping_pong_championship/referee/service"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/join", services.Join).Methods(http.MethodPost)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	errs := make(chan error)

	go func() {
		log.Println("Listening on localhost:8080")
		err := server.ListenAndServe()
		switch err {
		case http.ErrServerClosed:
		default:
			errs <- err
		}
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
		sig := <-c
		err := server.Shutdown(context.Background())
		errs <- fmt.Errorf("Signal: %s, ShutdownErr: %v", sig, err)
	}()

	err := <-errs

	log.Printf("terminated %s", err.Error())

}
