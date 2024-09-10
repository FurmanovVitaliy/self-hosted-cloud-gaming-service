package games

type Game struct {
	ID          string   `json:"id" bson:"_id,omitempty"`
	Name        string   `json:"name" bson:"name"`
	Path        string   `json:"-" bson:"path"`
	Url         string   `json:"url,omitempty" bson:"url"`
	Poster      string   `json:"poster,omitempty" bson:"logo"`
	Platform    string   `json:"platform,omitempty" bson:"platform"`
	Rating      float64  `json:"rating,omitempty" bson:"rating"`
	Summary     string   `json:"summary,omitempty" bson:"summary"`
	Videos      []string `json:"videos,omitempty" bson:"videos"`
	ReleaseDate int      `json:"release,omitempty" bson:"releaseDate"`
	IsGame      bool     `json:"-" bson:"isGame"`
}

// TODO: inject DTO model to create method
type CreateGameDTO struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Url    string `json:"url,omitempty"`
	Poster string `json:"poster,omitempty"`
}

type UpdateGameDTO struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
