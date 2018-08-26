package messages

type Message struct {
	Id      string `json:"_id"`
	RoomId  string `json:"rid"`
	Message string `json:"msg"`
}

const DefaultRoom = "GENERAL"

func NewSendMessageToDefaultRoom(message string) *MethodCall {
	return NewSendMessage(DefaultRoom, message)
}

func NewSendMessage(roomId, message string) *MethodCall {
	params := []interface{}{
		Message{
			Id:      genUniqueId(),
			RoomId:  roomId,
			Message: message,
		},
	}

	return NewMethodCall("sendMessage", params)
}
