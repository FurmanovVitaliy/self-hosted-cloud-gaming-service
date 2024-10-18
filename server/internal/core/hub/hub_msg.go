package hub

type ChatMsg struct {
	RoomUUID string `json:"room_uuid"`
	Username string `json:"username"`
	Content  string `json:"content"`
}

type WrtcMsg struct {
	ContentType string `json:"content_type"`
	Content     string `json:"content"`
}
