package server

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rainu/rocketchat-user-proxy/internal/client"
	"github.com/rainu/rocketchat-user-proxy/internal/server/handler"
	"net/http"
	"os"
)

func NewRouter(chat client.RocketChat) http.Handler {
	router := mux.NewRouter()

	// RESTful API
	router.Handle("/api/v1/send/u/{user}", handler.NewSendUserHandler(chat)).Methods(http.MethodPost)
	router.Handle("/api/v1/trigger/u/{user}", handler.NewTriggerUserHandler(chat)).Methods(http.MethodPost)
	router.Handle("/api/v1/send/r/{room}", handler.NewSendRoomHandler(chat)).Methods(http.MethodPost)
	router.Handle("/api/v1/trigger/r/{room}", handler.NewTriggerRoomHandler(chat)).Methods(http.MethodPost)

	return handlers.LoggingHandler(os.Stdout, router)
}
