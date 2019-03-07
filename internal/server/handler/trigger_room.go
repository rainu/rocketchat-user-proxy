package handler

import (
	"github.com/gorilla/mux"
	"github.com/rainu/rocketchat-user-proxy/internal/client"
	"io/ioutil"
	"net/http"
)

type triggerRoomHandler struct {
	Chat client.RocketChat
}

func NewTriggerRoomHandler(chat client.RocketChat) http.Handler {
	return &triggerRoomHandler{
		Chat: chat,
	}
}

func (s *triggerRoomHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	room := vars["room"]

	rawMessage, err := ioutil.ReadAll(request.Body)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	//send message
	s.Chat.TriggerRoom(string(rawMessage), room)

	writer.WriteHeader(http.StatusCreated)
}
