package handler

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"rocketchat-user-proxy/client"
)

type sendUserHandler struct {
	Chat client.RocketChat
}

func NewSendUserHandler(chat client.RocketChat) http.Handler {
	return &sendUserHandler{
		Chat: chat,
	}
}

func (s *sendUserHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	user := vars["user"]

	rawMessage, err := ioutil.ReadAll(request.Body)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	//send message
	s.Chat.SendDirectMessage(string(rawMessage), user)

	writer.WriteHeader(http.StatusCreated)
}
