package library

import (
	"log"

	"github.com/Henry-Sarabia/igdb/v2"
)

// TODO: ID AND TOKEN move to config
func getExtraInfoByName(name string) (fullName, url, img string, isGame bool) {
	client := igdb.NewClient(ID, TOKEN, nil)
	game, err := client.Games.Search(
		name,
		igdb.SetLimit(2),
		igdb.SetFields("cover", "name", "url"),
		igdb.SetFilter("rating", igdb.OpGreaterThan, "50"),
		igdb.SetFilter("cover", igdb.OpNotEquals, "null"),
	)
	if err != nil {
		log.Printf("In IGBD not found: %s", name)
		return name, "", "", false
	}
	cover, _ := client.Covers.Get(game[0].Cover, igdb.SetFields("image_id")) // retrieve cover IDs
	title, _ := cover.SizedURL(igdb.Size1080p, 1)                            // resize to largest image available
	return game[0].Name, game[0].URL, title, true
}
