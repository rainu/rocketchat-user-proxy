package messages

func NewDeleteMessage(messageId string) *MethodCall {
	params := []interface{}{
		map[string]string{"_id": messageId},
	}

	return NewMethodCall("deleteMessage", params)
}
