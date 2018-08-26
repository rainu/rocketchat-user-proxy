package server

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"rocketchat-user-proxy/client"
	"rocketchat-user-proxy/server/handler"
)

func NewRouter(chat client.RocketChat) *mux.Router {
	router := mux.NewRouter()

	// RESTful API
	router.Handle("/api/v1/send/u/{user}", handlers.LoggingHandler(os.Stdout, handler.NewSendUserHandler(chat))).
		Methods(http.MethodPost)
	router.Handle("/api/v1/send/r/{room}", handlers.LoggingHandler(os.Stdout, handler.NewSendRoomHandler(chat))).
		Methods(http.MethodPost)

	return router
}
