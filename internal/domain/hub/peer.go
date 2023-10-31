package hub

type Peer struct {
	ID       string
	Host     bool
	RoomID   string
	Username string
	Input    chan []byte
}
type GetPeersRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
