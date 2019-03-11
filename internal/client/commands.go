package client

import (
	"encoding/json"
	"github.com/rainu/rocketchat-user-proxy/internal/log"
	"github.com/rainu/rocketchat-user-proxy/internal/messages"
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
	roomId := rc.getRoomIdForUsername(recipient)

	if roomId != "" {
		//...and can send the message
		rc.sendMessage(messages.NewSendMessage(roomId, message))
	} else {
		log.Warning.Printf("The user '%v' is not known!\n", recipient)
	}
}

func (rc *rcClient) getRoomIdForUsername(recipient string) string {
	//do we know the user already?
	if _, alreadyKnown := rc.userDictionary[recipient]; !alreadyKnown {

		//we have to determine the roomId first...
		rawResp := rc.sendMessageAndWaitForResponse(messages.NewCreateDirectMessage(recipient))

		resp := messages.CreateDirectMessageResponse{}
		err := json.Unmarshal([]byte(rawResp), &resp)

		if err == nil {
			//now we know the roomId
			rc.userDictionary[recipient] = resp.Result.RoomId
		}
	}

	return rc.userDictionary[recipient]
}

func (rc *rcClient) SendRoomMessage(message string, rooms ...string) {
	for _, room := range rooms {
		rc.internalSendRoomMessage(message, room)
	}
}

func (rc *rcClient) internalSendRoomMessage(message string, room string) {
	roomId := rc.getRoomIdForRoomName(room)

	if roomId != "" {
		//...and can send the message
		rc.sendMessage(messages.NewSendMessage(roomId, message))
	} else {
		log.Warning.Printf("The room '%v' is not known!\n", room)
	}
}

func (rc *rcClient) getRoomIdForRoomName(room string) string {
	if _, alreadyKnown := rc.roomDictionary[room]; !alreadyKnown {

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
		}
	}

	return rc.roomDictionary[room]
}

func (rc *rcClient) DeleteMessage(messageId string) {
	rc.sendMessage(messages.NewDeleteMessage(messageId))
}

func (rc *rcClient) TriggerUser(triggerMessage string, recipients ...string) {
	for _, recipient := range recipients {
		rc.triggerInternalUser(triggerMessage, recipient)
	}
}

func (rc *rcClient) triggerInternalUser(triggerMessage, recipient string) {
	roomId := rc.getRoomIdForUsername(recipient)
	rc.sendAndDeleteMessage(triggerMessage, roomId)
}

func (rc *rcClient) TriggerRoom(triggerMessage string, rooms ...string) {
	for _, room := range rooms {
		rc.triggerInternalRoom(triggerMessage, room)
	}
}

func (rc *rcClient) triggerInternalRoom(triggerMessage, room string) {
	roomId := rc.getRoomIdForRoomName(room)
	rc.sendAndDeleteMessage(triggerMessage, roomId)
}

func (rc *rcClient) sendAndDeleteMessage(triggerMessage, roomId string) {
	rawResp := rc.sendMessageAndWaitForResponse(messages.NewSendMessage(roomId, triggerMessage))

	resp := messages.SendMessageResponse{}
	err := json.Unmarshal([]byte(rawResp), &resp)

	if err == nil {
		rc.DeleteMessage(resp.Result.Id)
	}
}
