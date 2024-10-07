package scanner

type Game struct {
	Name     string
	Path     string
	Platform string
}

type HashRecord struct {
	ID   string `bson:"_id"`
	Hash string `bson:"hash"`
}
