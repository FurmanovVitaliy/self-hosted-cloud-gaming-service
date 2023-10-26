package games

type Game struct {
	ID     string `json:"id" bson:"_id,omitempty"`
	Name   string `json:"name" bson:"name"`
	Path   string `json:"-" bson:"path"`
	Url    string `json:"url,omitempty" bson:"url"`
	Logo   string `json:"logo,omitempty" bson:"logo"`
	IsGame bool   `json:"-" bson:"isGame"`
}

// TODO: inject DTO model to create method
type CreateGameDTO struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Url  string `json:"url,omitempty"`
	Logo string `json:"logo,omitempty"`
}

type UpdateGameDTO struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
