package handler

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"rocketchat-user-proxy/client"
)

type triggerUserHandler struct {
	Chat client.RocketChat
}

func NewTriggerUserHandler(chat client.RocketChat) http.Handler {
	return &triggerUserHandler{
		Chat: chat,
	}
}

func (s *triggerUserHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	user := vars["user"]

	rawMessage, err := ioutil.ReadAll(request.Body)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	//send message
	s.Chat.TriggerUser(string(rawMessage), user)

	writer.WriteHeader(http.StatusCreated)
}
