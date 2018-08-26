package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"rocketchat-user-proxy/client"
	"rocketchat-user-proxy/log"
	"rocketchat-user-proxy/server"
	"syscall"
	"time"
)

func main() {
	rc := client.NewRocketChat("url")
	err := rc.Start()

	if err != nil {
		log.Error.Printf("Could not establish connection. Error: %v\n", err)
		return
	}

	rc.LoginWithPassword("user", "password")

	httpServer := startServer(rc)

	//wait for interruption
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	httpServer.Shutdown(ctx)
	rc.Stop()
}

func startServer(chat client.RocketChat) *http.Server {
	router := server.NewRouter(chat)
	httpServer := &http.Server{Addr: fmt.Sprintf(":%v", 8080), Handler: router}

	go func() {
		httpServer.ListenAndServe()
	}()

	return httpServer
}
