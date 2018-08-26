package client

import (
	"encoding/json"
	"rocketchat-user-proxy/log"
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

func (rc *rcClient) SendDirectMessage(message string, recipients ...string) {
	if len(recipients) == 0 {
		rc.sendMessage(messages.NewSendMessageToDefaultRoom(message))
	} else {
		for _, recipient := range recipients {
			rc.internalSendDirectMessage(message, recipient)
		}
	}
}

func (rc *rcClient) internalSendDirectMessage(message string, recipient string) {
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

func (rc *rcClient) SendRoomMessage(message string, rooms ...string) {
	for _, room := range rooms {
		rc.internalSendRoomMessage(message, room)
	}
}

func (rc *rcClient) internalSendRoomMessage(message string, room string) {
	if roomId, alreadyKnown := rc.roomDictionary[room]; alreadyKnown {
		rc.sendMessage(messages.NewSendMessage(roomId, message))
		return
	}

	//we have to determine the roomId first
	rawResp := rc.sendMessageAndWaitForResponse(messages.NewGetRooms())

	resp := messages.GetRoomsResponse{}
	err := json.Unmarshal([]byte(rawResp), &resp)

	if err == nil {
		//now we know a lot of roomIds
		for _, updateEntry := range resp.Result.Update {
			if updateEntry["t"] == "c" {
				//this is a channelEntry
				rc.roomDictionary[updateEntry["name"].(string)] = updateEntry["_id"].(string)
			}
		}

		if roomId, isNowKnown := rc.roomDictionary[room]; isNowKnown {
			//...and can send the message
			rc.sendMessage(messages.NewSendMessage(roomId, message))
		} else {
			log.Warning.Printf("The room '%v' is not known!\n", room)
		}
	}
}
