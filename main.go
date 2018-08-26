package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"rocketchat-user-proxy/client"
	"rocketchat-user-proxy/config"
	"rocketchat-user-proxy/log"
	"rocketchat-user-proxy/server"
	"syscall"
	"time"
)

func main() {
	cfg := config.New()

	rc := client.NewRocketChat(cfg.Url)
	err := rc.Start()

	if err != nil {
		log.Error.Printf("Could not establish connection. Error: %v\n", err)
		return
	}

	if cfg.PasswordHash != "" {
		rc.LoginWithHash(cfg.Username, cfg.PasswordHash)
	} else {
		rc.LoginWithPassword(cfg.Username, cfg.PasswordPlain)
	}

	httpServer := startServer(rc, cfg)

	//wait for interruption
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	httpServer.Shutdown(ctx)
	rc.Stop()
}

func startServer(chat client.RocketChat, cfg *config.Config) *http.Server {
	router := server.NewRouter(chat)
	httpServer := &http.Server{Addr: fmt.Sprintf(":%v", cfg.BindPort), Handler: router}

	go func() {
		httpServer.ListenAndServe()
	}()

	return httpServer
}
