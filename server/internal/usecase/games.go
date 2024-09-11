package usecase

import (
	"context"
	"net/http"
	"time"

	"github.com/FurmanovVitaliy/pixel-cloud/pkg/errors"
)

// dto
type GetGamesReq struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Url         string   `json:"url,omitempty" `
	Poster      string   `json:"poster,omitempty"`
	Platform    string   `json:"platform,omitempty"`
	Rating      float64  `json:"rating,omitempty" `
	Summary     string   `json:"summary,omitempty"`
	Videos      []string `json:"videos,omitempty" `
	ReleaseDate int      `json:"release,omitempty"`
}

// custom errors
var (
	ErrFetchGames      = errors.New(http.StatusInternalServerError, "GS", "000001", "error while fetching games")
	ErrGameUnavailable = errors.New(http.StatusNotFound, "GS", "000002", "game not available")
)

func (u *UseCase) GetGames() (games []GetGamesReq, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	g, err := u.gameService.GetAll(ctx)
	if err != nil {
		return []GetGamesReq{}, ErrFetchGames
	}
	for _, game := range g {
		games = append(games, GetGamesReq{
			ID:          game.ID,
			Name:        game.Name,
			Url:         game.Url,
			Poster:      game.Poster,
			Platform:    game.Platform,
			Rating:      game.Rating,
			Summary:     game.Summary,
			Videos:      game.Videos,
			ReleaseDate: game.ReleaseDate,
		})
	}
	return games, nil
}
