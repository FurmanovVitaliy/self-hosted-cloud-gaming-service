package library

import (
	"cloud/internal/domain/games"
	"cloud/pkg/logger"
	hashsum "cloud/pkg/utils/heshsum"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Henry-Sarabia/igdb/v2"
	"github.com/sahilm/fuzzy"
)

type gameCandidate struct {
	name       string
	extantion  string
	path       string
	matchScore int
}

type gameDir struct {
	platform   string
	name       string
	path       string
	candidates []gameCandidate
	oneMatch   bool
}

type gameSearch struct {
	extantions     []string
	directories    []string
	namesToCompare []string //for easier parsing while searching
	gameDirs       []gameDir
	igdb           *igdb.Client
	logger         *logger.Logger
}

func NewSerchEngine(extantions, directories, names []string, id, token string, logger *logger.Logger) *gameSearch {
	igdbClient := igdb.NewClient(id, token, nil)
	if id == "" || token == "" {
		logger.Error("IGDB credentials are not set")
	}
	return &gameSearch{
		extantions:     extantions,
		directories:    directories,
		namesToCompare: names,
		igdb:           igdbClient,
		logger:         logger,
	}
}

func (e *gameSearch) ScanLibrary() ([]games.Game, string, error) {
	e.logger.Info("scaning is started")
	for _, path := range e.directories {
		e.logger.Info("scanning directory with path: ", path)
		entries, err := os.ReadDir(path)
		if err != nil {
			e.logger.Error(err)
			return nil, "", err
		}
		for _, entry := range entries {
			if entry.IsDir() {
				if e.notSystemDirectory(entry.Name()) {
					e.gameDirs = append(e.gameDirs, gameDir{name: entry.Name(), path: path + "/" + entry.Name()})
				}
			}
		}
	}
	if len(e.gameDirs) == 0 {
		return nil, "", fmt.Errorf("no game directories found")
	}

	for i, gDir := range e.gameDirs {
		err := filepath.Walk(gDir.path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				e.logger.Error(err)
				return err
			}
			if !info.IsDir() {
				for _, ext := range e.extantions {
					if strings.HasSuffix(info.Name(), ext) {
						e.gameDirs[i].candidates = append(e.gameDirs[i].candidates, gameCandidate{name: strings.TrimSuffix(info.Name(), "."+ext), extantion: ext, path: path})
					}
				}
			}
			return nil
		})
		if err != nil {
			return nil, "", err
		}
	}
	var g []games.Game

	for _, gDir := range e.gameDirs {
		e.choseCandidate(&gDir)
		g = append(g, games.Game{
			Name:     e.namePrettyfier(gDir.name),
			Path:     gDir.candidates[0].path,
			Platform: gDir.platform,
		})
	}

	e.logger.Info("scaning is finished")
	hash := hashsum.Hashsum(g)
	return g, hash, nil
}

func (e *gameSearch) GetInfoFromIGDB(gms []games.Game) ([]games.Game, error) {
	e.logger.Info("getting exra info from IGDB")
	for i, game := range gms {
		info, err := e.igdb.Games.Search(
			game.Name,
			igdb.SetFields("cover", "name", "url", "total_rating", "summary", "videos", "first_release_date"),
			igdb.SetFilter("cover", igdb.OpNotEquals, "null"),
			igdb.SetFilter("version_parent", igdb.OpEquals, "null"),
			igdb.SetFilter("rating", igdb.OpGreaterThan, "20"),
			igdb.SetFilter("external_games", igdb.OpNotEquals, "null"),

			igdb.SetLimit(1),
		)
		if err != nil {
			e.logger.Error(err)
			e.logger.Warnf("game %s IGBD not found", game.Name)
			continue
		}

		cover, _ := e.igdb.Covers.Get(info[0].Cover, igdb.SetFields("image_id"))
		poster, _ := cover.SizedURL(igdb.Size1080p, 1)
		videosStringURL := []string{}

		if len(info[0].Videos) != 0 {
			for _, v := range info[0].Videos {
				video, err := e.igdb.GameVideos.Get(v, igdb.SetFields("name", "video_id"))
				if err != nil {
					e.logger.Error(err)
					continue
				}
				if video.Name == "Trailer" {
					stringURl := "https://www.youtube.com/watch?v=" + video.VideoID
					videosStringURL = append(videosStringURL, stringURl)

				}
			}
		}

		release := time.Unix(int64(info[0].FirstReleaseDate), 0).Year()

		(gms)[i].Name = info[0].Name
		(gms)[i].Url = info[0].URL
		(gms)[i].Rating = info[0].TotalRating
		(gms)[i].Summary = info[0].Summary
		(gms)[i].Videos = videosStringURL
		(gms)[i].Poster = poster
		(gms)[i].ReleaseDate = release
		(gms)[i].IsGame = true
	}
	return gms, nil
}

func (e *gameSearch) choseCandidate(d *gameDir) {
	var maxScore int
	var maxScoreIndex int
	e.sizeScore(d)
	e.nameScore(d)
	e.exeScore(d)
	if d.oneMatch {
		return
	}
	if len(d.candidates) == 0 {
		return
	}
	for i, c := range d.candidates {
		if c.matchScore > maxScore {
			maxScore = c.matchScore
			maxScoreIndex = i
		}
	}
	d.candidates = []gameCandidate{d.candidates[maxScoreIndex]}
	switch d.candidates[0].extantion {
	case "exe":
		d.platform = "windows"
	case "nsp":
		d.platform = "switch"
	case "x86_64":
		d.platform = "linux"
	}
}

func (e *gameSearch) sizeScore(d *gameDir) {
	var maxFileSize int64
	if len(d.candidates) == 1 {
		d.oneMatch = true
		return
	}
	if len(d.candidates) == 0 {
		return
	}

	for i, c := range d.candidates {
		fileInfo, err := os.Stat(c.path)
		if err != nil {
			e.logger.Error(err)
			continue
		}
		if fileInfo.Size() > maxFileSize {
			maxFileSize = fileInfo.Size()
			d.candidates[i].matchScore += 4
		}
	}
}
func (e *gameSearch) exeScore(d *gameDir) {
	if len(d.candidates) == 1 {
		d.oneMatch = true
		return
	}
	if len(d.candidates) == 0 {
		return
	}
	for i, c := range d.candidates {
		if c.extantion == "exe" {
			_, err := os.Stat(c.path + ".ppdb")
			if err == nil {
				d.candidates[i].matchScore += 4
			} else {
				continue
			}
		}
	}
}
func (e *gameSearch) nameScore(d *gameDir) {
	if len(d.candidates) == 1 {
		d.oneMatch = true
		return
	}
	if len(d.candidates) == 0 {
		return
	}
	for i, c := range d.candidates {
		c.name = e.namePrettyfier(c.name)

		matches := fuzzy.Find(c.name, e.namesToCompare)
		if len(matches) > 0 {
			c.name = e.namesToCompare[matches[0].Index]
			d.candidates[i].matchScore += 8
		} else if c.name == "launcher" {
			d.candidates[i].matchScore += 10
		} else {
			levMatch := e.levenshteinDistanceInPersent(c.name, e.namePrettyfier(d.name))
			switch {
			case levMatch > 80:
				d.candidates[i].matchScore += 8
			case levMatch > 60:
				d.candidates[i].matchScore += 6
			case levMatch > 40:
				d.candidates[i].matchScore += 4
			case levMatch > 20:
				d.candidates[i].matchScore += 2
			}
		}
	}
}

func (e *gameSearch) namePrettyfier(name string) string {
	wordsToRemove := []string{"dlc", "crack", "patch", "update", "fix", "repack", "multi", "rus", "eng", "engrus", "engrusmulti", "multie", "pack"}
	regexPattern := "\\b(" + strings.Join(wordsToRemove, "|") + ")\\b"

	prettyName := regexp.MustCompile(`\[[^\[\]]*\]`).ReplaceAllString(name, "")
	prettyName = regexp.MustCompile(`\([^()]*\)`).ReplaceAllString(prettyName, "")
	prettyName = regexp.MustCompile(`([^ ]) ([A-Z])`).ReplaceAllString(prettyName, "$1 $2")
	prettyName = regexp.MustCompile(`[^a-zA-Z0-9'"-: ]+`).ReplaceAllString(prettyName, "")
	prettyName = strings.ToLower(prettyName)
	prettyName = regexp.MustCompile(regexPattern).ReplaceAllString(prettyName, "")
	prettyName = regexp.MustCompile(`\s+`).ReplaceAllString(prettyName, " ")

	return prettyName
}

func (e *gameSearch) levenshteinDistanceInPersent(s1, s2 string) int {
	lenS1 := len(s1)
	lenS2 := len(s2)

	matrix := make([][]int, lenS1+1)
	for i := range matrix {
		matrix[i] = make([]int, lenS2+1)
	}

	for i := 0; i <= lenS1; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= lenS2; j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= lenS1; i++ {
		for j := 1; j <= lenS2; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}
			matrix[i][j] = min(matrix[i-1][j]+1, matrix[i][j-1]+1, matrix[i-1][j-1]+cost)
		}
	}

	return 100 - ((matrix[lenS1][lenS2]) * 100 / 21)
}

func (e *gameSearch) notSystemDirectory(name string) bool {
	systemDirectories := []string{"steam", "proton", "wine", "yuzu", "cemu", "rpcs3"}
	name = e.namePrettyfier(name)
	bool := true

	for _, dir := range systemDirectories {
		if strings.Contains(name, dir) {
			bool = false
		}
	}
	return bool
}
