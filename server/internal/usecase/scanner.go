package usecase

import (
	"context"

	"github.com/FurmanovVitaliy/pixel-cloud/internal/domain/game"
)

func (u *UseCase) ScanLibrary() {
	//checking for changes in the library
	if !u.scannerSerrvice.CheckForChanges() {
		return
	}
	//clear the database
	err := u.gameService.DeleteAll(context.Background())
	if err != nil {
		return
	}
	//scan the library
	candidates, err := u.scannerSerrvice.Scan()
	if err != nil {
		return
	}
	//fetch additional information about the games
	for _, candidate := range candidates {
		var g game.Game
		extraInfo, err := u.gameInfoService.GetExtraInfoByName(candidate.Name)
		if err != nil {
			continue
		}
		//create a new game
		g = game.Game{
			Name:        extraInfo.Name,
			Path:        candidate.Path,
			Url:         extraInfo.Url,
			Poster:      extraInfo.Poster,
			Platform:    candidate.Platform,
			Rating:      extraInfo.Rating,
			Summary:     extraInfo.Summary,
			Videos:      extraInfo.Videos,
			ReleaseDate: extraInfo.ReleaseDate,
			IsGame:      true,
		}
		//create the game in the database
		_, err = u.gameService.Create(context.Background(), g)
		if err != nil {
			continue
		}
	}
}
