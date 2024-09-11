package igdb

import (
	"fmt"
	"time"

	"github.com/FurmanovVitaliy/pixel-cloud/pkg/logger"

	"github.com/Henry-Sarabia/igdb/v2"
)

type GameExtraInfo struct {
	Name        string
	Url         string
	Poster      string
	Rating      float64
	Summary     string
	Videos      []string
	ReleaseDate int
}

type extraDB struct {
	client *igdb.Client
	logger *logger.Logger
}

func NewClient(id, token string) (client *igdb.Client, err error) {
	client = igdb.NewClient(id, token, nil)
	_, err = client.Games.Search("zelda")
	if err != nil {
		return nil, err
	}
	return

}

func New(client *igdb.Client, logger *logger.Logger) *extraDB {
	return &extraDB{
		client: client,
		logger: logger,
	}
}
func (d *extraDB) GetExtraInfoByName(name string) (gameInfo GameExtraInfo, err error) {
	info, err := d.client.Games.Search(
		name,
		igdb.SetFields("cover", "name", "url", "total_rating", "summary", "videos", "first_release_date"),
		igdb.SetFilter("cover", igdb.OpNotEquals, "null"),
		igdb.SetFilter("version_parent", igdb.OpEquals, "null"),
		igdb.SetFilter("rating", igdb.OpGreaterThan, "20"),
		igdb.SetFilter("external_games", igdb.OpNotEquals, "null"),
		igdb.SetLimit(1),
	)

	if err != nil {
		d.logger.Warnf("Error getting game info for %s: %v", name, err)
		return GameExtraInfo{
			Name: name,
		}, err
	}

	if len(info) == 0 {
		return gameInfo, fmt.Errorf("game info is not found")
	}

	cover, _ := d.client.Covers.Get(info[0].Cover, igdb.SetFields("image_id"))
	poster, _ := cover.SizedURL(igdb.Size1080p, 1)
	videosStringURL := []string{}

	if len(info[0].Videos) != 0 {
		for _, v := range info[0].Videos {
			video, err := d.client.GameVideos.Get(v, igdb.SetFields("name", "video_id"))
			if err != nil {
				// Handle error and continue with next video
				continue
			}
			if video.Name == "Trailer" {
				stringURL := "https://www.youtube.com/watch?v=" + video.VideoID
				videosStringURL = append(videosStringURL, stringURL)
			}
		}
	}

	release := time.Unix(int64(info[0].FirstReleaseDate), 0).Year()

	return GameExtraInfo{
		Name:        info[0].Name,
		Url:         info[0].URL,
		Poster:      poster,
		Rating:      info[0].TotalRating,
		Summary:     info[0].Summary,
		Videos:      videosStringURL,
		ReleaseDate: release,
	}, nil
}
