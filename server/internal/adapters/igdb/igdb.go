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
	Genres      []string
	Summary     string
	Videos      []string
	Images      []string
	ReleaseDate int
	AgeRating   string
	Developer   string
	Publisher   string
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
		igdb.SetFields("cover", "name", "url", "total_rating", "summary", "videos", "first_release_date", "screenshots", "artworks", "genres", "age_ratings", "involved_companies", "age_ratings"),
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

	videosURLs := []string{}

	if len(info[0].Videos) != 0 {
		for _, v := range info[0].Videos {
			video, err := d.client.GameVideos.Get(v, igdb.SetFields("name", "video_id"))
			if err != nil {
				// Handle error and continue with next video
				continue
			}
			if video.Name == "Trailer" {
				stringURL := "https://www.youtube.com/watch?v=" + video.VideoID
				videosURLs = append(videosURLs, stringURL)
			}
		}
	}

	screenshotsURLs := []string{}

	if len(info[0].Screenshots) != 0 {
		for _, s := range info[0].Screenshots {
			screenshot, err := d.client.Screenshots.Get(s, igdb.SetFields("url"))
			if err != nil {
				// Handle error and continue with next screenshot
				continue
			}
			screenshotsURLs = append(screenshotsURLs, screenshot.URL)
		}
	}

	if len(info[0].Artworks) != 0 {
		for _, a := range info[0].Artworks {
			artwork, err := d.client.Artworks.Get(a, igdb.SetFields("url"))
			if err != nil {
				// Handle error and continue with next artwork
				continue
			}
			screenshotsURLs = append(screenshotsURLs, artwork.URL)
		}
	}

	genres := []string{}

	if len(info[0].Genres) != 0 {
		for _, g := range info[0].Genres {
			genre, err := d.client.Genres.Get(g, igdb.SetFields("name"))
			if err != nil {
				// Handle error and continue with next genre
				continue
			}
			genres = append(genres, genre.Name)
		}
	}

	var developer string
	var publisher string

	if len(info[0].InvolvedCompanies) != 0 {
		for _, i := range info[0].InvolvedCompanies {
			company, err := d.client.InvolvedCompanies.Get(i, igdb.SetFields("company", "developer", "publisher"))
			if err != nil {
				// Handle error and continue with next company
				continue
			}

			if company.Developer {
				company1, err := d.client.Companies.Get(company.Company, igdb.SetFields("name"))
				if err != nil {
					// Handle error and continue with next company
					continue
				}
				developer = company1.Name
			}

			if company.Publisher {
				company1, err := d.client.Companies.Get(company.Company, igdb.SetFields("name"))
				if err != nil {
					// Handle error and continue with next company
					continue
				}
				publisher = company1.Name
			}
		}
	}

	var ageRating string

	if len(info[0].AgeRatings) != 0 {
		for _, a := range info[0].AgeRatings {
			rating, err := d.client.AgeRatings.Get(a, igdb.SetFields("rating", "category"))
			if err != nil {
				// Handle error and continue with next age rating
				continue
			}
			if rating.Category == 2 { //2 for PEGI rating
				switch rating.Rating {
				case 1:
					ageRating = "3+"
				case 2:
					ageRating = "7+"
				case 3:
					ageRating = "12+"
				case 4:
					ageRating = "16+"
				case 5:
					ageRating = "18+"
				default:
					ageRating = "0+"
				}
			}
		}
	}

	release := time.Unix(int64(info[0].FirstReleaseDate), 0).Year()

	return GameExtraInfo{
		Name:        info[0].Name,
		Url:         info[0].URL,
		Rating:      info[0].TotalRating,
		Summary:     info[0].Summary,
		Poster:      poster,
		Videos:      videosURLs,
		Images:      screenshotsURLs,
		ReleaseDate: release,
		Genres:      genres,
		AgeRating:   ageRating,
		Developer:   developer,
		Publisher:   publisher,
	}, nil
}
