package messages

type GetRoomsResponseResult struct {
	Update []map[string]interface{} `json:"update"`
	Remove []map[string]interface{} `json:"remove"`
}

type GetRoomsResponse struct {
	Id     string                 `json:"id"`
	Msg    string                 `json:"msg"`
	Result GetRoomsResponseResult `json:"result"`
}

func NewGetRooms() *MethodCall {
	params := []interface{}{
		map[string]int{"$date": 0},
	}

	return NewMethodCall("rooms/get", params)
}
