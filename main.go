package main

import (
	"os"
	"os/signal"
	"rocketchat-user-proxy/client"
	"rocketchat-user-proxy/log"
	"rocketchat-user-proxy/messages"
	"syscall"
)

func main() {
	rc := client.NewRocketChat("url")
	in, _, err := rc.Start()

	if err != nil {
		log.Error.Printf("Could not establish connection. Error: %v\n", err)
		return
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	in <- messages.NewLoginPlain("USERNAM", "PASSWORD")
	in <- messages.NewSendMessageToDefaultRoom("Erste progammierter Test4!")

	//wait for interruption
	<-stop
	rc.Stop()
}
