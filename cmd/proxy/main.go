package main

import (
	"context"
	"fmt"
	"github.com/rainu/rocketchat-user-proxy/internal/client"
	"github.com/rainu/rocketchat-user-proxy/internal/log"
	"github.com/rainu/rocketchat-user-proxy/internal/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := NewConfig()

	rc := client.NewRocketChat(cfg.WSUrl)
	err := rc.Start()

	if err != nil {
		log.Error.Printf("Could not establish connection. Error: %v\n", err)
		return
	}
	defer rc.Stop()

	if cfg.PasswordHash != "" {
		rc.LoginWithHash(cfg.Username, cfg.PasswordHash)
	} else {
		rc.LoginWithPassword(cfg.Username, cfg.PasswordPlain)
	}
	defer rc.Logout()

	httpServer := startServer(rc, cfg)
	defer func() {
		//gracefully shutdown the httpServer
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		httpServer.Shutdown(ctx)
	}()

	//wait for interruption
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}

func startServer(chat client.RocketChat, cfg *Config) *http.Server {
	router := server.NewRouter(chat)
	httpServer := &http.Server{Addr: fmt.Sprintf(":%v", cfg.BindPort), Handler: router}

	go func() {
		httpServer.ListenAndServe()
	}()

	return httpServer
}
