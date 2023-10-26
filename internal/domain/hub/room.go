package hub

type Room struct {
	ID    string
	Name  string
	Peers map[string]*Peer
}

type CreateRoomRequest struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type GetRoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
