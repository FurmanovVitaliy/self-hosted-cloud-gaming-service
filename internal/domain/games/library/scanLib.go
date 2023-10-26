package library

import (
	"cloud/internal/domain/games"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sahilm/fuzzy"
)

func regExFilter(name string) string {
	// Создаем регулярное выражение для поиска заглавных букв и добавления перед ними пробела
	re := regexp.MustCompile(`([A-Z])`)
	result := re.ReplaceAllString(name, " $1")
	// Удаляем начальный пробел, если он есть
	result = strings.TrimSpace(result)
	// Удаляем все специальные символы и цифры
	result = regexp.MustCompile(`[^a-zA-Z ]+`).ReplaceAllString(result, "")
	// Удаляем приставку "exeppdb" с конца строки
	result = strings.TrimSuffix(result, "exeppdb")
	return result
}

func locDbFilter(name string) string {
	matches := fuzzy.Find(name, locGamesDatabase)
	if len(matches) > 0 {
		name = locGamesDatabase[matches[0].Index]
	}
	return name
}

func getGamePath(path string) string {
	path = strings.TrimSuffix(path, ".ppdb")
	return path
}
func ScanLib() ([]games.Game, error) {
	var gameLib []games.Game

	err := filepath.Walk(DIRECTORY_TO_SEARCH, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("Error occured")
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == FILE_EXTANTION {
			name, url, logo, isGame := getExtraInfoByName(locDbFilter(regExFilter(info.Name())))
			gameCandidate := games.Game{Path: getGamePath(path), Name: name, Url: url, Logo: logo, IsGame: isGame}
			if isGame {
				gameLib = append(gameLib, gameCandidate)
			}
		}
		return nil
	})
	if err != nil {
		log.Println("Error occured")
		return nil, err
	}
	return gameLib, nil
}
