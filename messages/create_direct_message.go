package messages

type CreateDirectMessageResponseResult struct {
	RoomId string `json:"rid"`
}

type CreateDirectMessageResponse struct {
	Id     string                            `json:"id"`
	Msg    string                            `json:"msg"`
	Result CreateDirectMessageResponseResult `json:"result"`
}

func NewCreateDirectMessage(targetUsername string) *MethodCall {
	params := []interface{}{
		targetUsername,
	}

	return NewMethodCall("createDirectMessage", params)
}
