package client

import (
	"encoding/json"
	"github.com/rainu/rocketchat-user-proxy/internal/messages"
)

type History interface {
	AddIncomingMessage(string)
	AddOutgoingMessage(*messages.MethodCall)

	WaitForResponse(*messages.MethodCall) string
}

type history struct {
	in  map[string]string
	out map[string]*messages.MethodCall

	waits map[string]chan string
}

func NewHistory() History {
	return &history{
		in:    make(map[string]string),
		out:   make(map[string]*messages.MethodCall),
		waits: make(map[string]chan string),
	}
}

func (h *history) AddIncomingMessage(message string) {
	//try to convert in GeneralMessage
	genResp := &messages.GeneralResponse{}
	err := json.Unmarshal([]byte(message), genResp)

	if err == nil && genResp.Msg == "result" {
		//we only want to store relational messages
		h.in[genResp.Id] = message

		if waitChan, contains := h.waits[genResp.Id]; contains {
			waitChan <- message //send signal
		}
	}
}

func (h *history) AddOutgoingMessage(message *messages.MethodCall) {
	h.out[message.Id] = message
}

func (h *history) WaitForResponse(message *messages.MethodCall) string {
	if response, contains := h.in[message.Id]; contains {
		//we have already receive the response
		return response
	} else {
		//we have to wait for the response

		//register wait channel
		h.waits[message.Id] = make(chan string, 1)

		//wait for signal
		return <-h.waits[message.Id]
	}
}
