package player

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"ping_pong_championship/player/config"
	"ping_pong_championship/player/transport"
	"syscall"
)

func main() {

	config.ParseFlags()

	server := config.GetHttpServer()
	server.Handler = transport.MakeHandler()

	errs := make(chan error)

	go func() {
		log.Println("Listening on localhost:8081")
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
		errs <- fmt.Errorf("Singal: %s, ShutdownErr: %v", sig, err)
	}()

	err := <-errs

	log.Printf("terminated %s", err.Error())
}
