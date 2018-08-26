package client

import (
	"encoding/json"
	"rocketchat-user-proxy/messages"
)

func (rc *rcClient) sendMessage(message *messages.MethodCall) {
	rc.chanIn <- message
	rc.history.AddOutgoingMessage(message)
}

func (rc *rcClient) sendMessageAndWaitForResponse(message *messages.MethodCall) string {
	rc.sendMessage(message)

	return rc.history.WaitForResponse(message)
}

func (rc *rcClient) LoginWithPassword(username, password string) {
	rc.sendMessageAndWaitForResponse(messages.NewLoginPlain(username, password))
}

func (rc *rcClient) LoginWithHash(username, passwordHash string) {
	rc.sendMessageAndWaitForResponse(messages.NewLoginHash(username, passwordHash))
}

func (rc *rcClient) Logout() {
	rc.sendMessage(messages.NewLogout())
}

func (rc *rcClient) SendMessage(message string, recipients ...string) {
	if len(recipients) == 0 {
		rc.sendMessage(messages.NewSendMessageToDefaultRoom(message))
	} else {
		for _, recipient := range recipients {
			rc.sendDirectMessage(message, recipient)
		}
	}
}

func (rc *rcClient) sendDirectMessage(message string, recipient string) {
	//do we know the user already?
	if roomId, alreadyKnown := rc.userDictionary[recipient]; alreadyKnown {
		rc.sendMessage(messages.NewSendMessage(roomId, message))
		return
	}

	//we have to determine the roomId first...
	rawResp := rc.sendMessageAndWaitForResponse(messages.NewCreateDirectMessage(recipient))

	resp := messages.CreateDirectMessageResponse{}
	err := json.Unmarshal([]byte(rawResp), &resp)

	if err == nil {
		//now we know the roomId
		rc.userDictionary[recipient] = resp.Result.RoomId

		//...and can send the message
		rc.sendMessage(messages.NewSendMessage(resp.Result.RoomId, message))
	}
}
