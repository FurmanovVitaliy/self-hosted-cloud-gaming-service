package messages

import "encoding/json"

const (
	RTC_OFFER               = "offer"
	RTC_ANSWER              = "answ"
	RTC_SERVER_CANDIDATE    = "sc"
	RTC_CLIENT_CANDIDATE    = "cc"
	RTC_SIGNAL              = "sig"
	RtcIceGatheringComplete = "serverIceGatheringComplete"
	RtcConnectionReady      = "connectionReady"
)

var (
	RtcConnectionClosed = NewMessage("", "", RTC_SIGNAL, "connectionClosed")
)

type Message struct {
	RoomUUID    string `json:"room_uuid"`
	Username    string `json:"username"`
	ContentType string `json:"content_type"`
	Content     string `json:"content"`
}

func NewMessage(roomUUID, username, contentType, content string) *Message {
	return &Message{
		RoomUUID:    roomUUID,
		Username:    username,
		ContentType: contentType,
		Content:     content,
	}
}

func (m *Message) Marshal() []byte {
	marshal, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return marshal
}
