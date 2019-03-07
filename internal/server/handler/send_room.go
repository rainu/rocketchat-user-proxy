package handler

import (
	"github.com/gorilla/mux"
	"github.com/rainu/rocketchat-user-proxy/internal/client"
	"io/ioutil"
	"net/http"
)

type sendRoomHandler struct {
	Chat client.RocketChat
}

func NewSendRoomHandler(chat client.RocketChat) http.Handler {
	return &sendRoomHandler{
		Chat: chat,
	}
}

func (s *sendRoomHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	room := vars["room"]

	rawMessage, err := ioutil.ReadAll(request.Body)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	//send message
	s.Chat.SendRoomMessage(string(rawMessage), room)

	writer.WriteHeader(http.StatusCreated)
}
