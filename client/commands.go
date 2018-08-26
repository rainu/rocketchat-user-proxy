package client

import "rocketchat-user-proxy/messages"

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
		//TODO
	}
}
