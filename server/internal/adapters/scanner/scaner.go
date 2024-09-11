package scanner

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/FurmanovVitaliy/pixel-cloud/pkg/logger"
)

type hashStorage interface {
	GetHash(ctx context.Context, id string) (string, error)
	UpsertHash(ctx context.Context, id string, hash string) error
}

type scaner struct {
	ctx     context.Context
	logger  *logger.Logger
	params  *params
	storage hashStorage
}

type gameCandidate struct {
	name       string
	path       string
	extantion  string
	matchScore int
}

type gameDir struct {
	name          string
	path          string
	Candidates    []gameCandidate
	bestCandidate gameCandidate
}

func New(ctx context.Context, logger *logger.Logger, params *params, storage hashStorage) *scaner {
	return &scaner{
		ctx:     ctx,
		logger:  logger,
		params:  params,
		storage: storage,
	}
}

func (s *scaner) CheckForChanges() bool {
	newHash, err := getFoldersHash(s.params.targetDirs)
	if err != nil {
		s.logger.Error(err)
		return true
	}

	prevHash, err := s.storage.GetHash(context.Background(), "hash")
	if err != nil {
		s.logger.Warn(err)
		err = s.storage.UpsertHash(context.Background(), "hash", newHash)
		if err != nil {
			s.logger.Error(err)
		}
	}

	if newHash == prevHash {
		s.logger.Info("no changes in library")
		return false
	}

	err = s.storage.UpsertHash(s.ctx, "hash", newHash)
	if err != nil {
		s.logger.Error(err)
	}
	return true
}

func (s *scaner) Scan() (games []Game, err error) {
	var wg sync.WaitGroup
	s.logger.Info("Scanning game library")
	gameDirs, _ := s.findGameDirs(s.params.targetDirs)
	for i := range gameDirs {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.logger.Debugf("Scanning %s for a launch file", gameDirs[i].path)
			s.findCandidates(&gameDirs[i])
			s.rateCandidates(&gameDirs[i])
			s.choseBestCandidate(&gameDirs[i])
			if gameDirs[i].bestCandidate.name != "" && gameDirs[i].bestCandidate.matchScore > 0 {
				games = append(games, Game{
					Name:     gameDirs[i].bestCandidate.name,
					Path:     gameDirs[i].bestCandidate.path,
					Platform: platformByExtantion(gameDirs[i].bestCandidate.extantion),
				})
			}
		}(i)
	}
	wg.Wait()
	return games, nil
}

func (s *scaner) findGameDirs(paths []string) ([]gameDir, error) {
	var d []gameDir
	for _, path := range paths {
		s.logger.Debugf("Scanning %s for game directories", path)
		entries, err := os.ReadDir(path)
		if err != nil {
			s.logger.Warn(err)
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				if notSystemDirectory(entry.Name(), s.params.excludeDirs) {
					d = append(d, gameDir{name: entry.Name(), path: path + "/" + entry.Name()})
				}
			}
		}
	}
	return d, nil
}

func (s *scaner) findCandidates(dir *gameDir) error {
	extSet := make(map[string]struct{})
	for _, ext := range s.params.targetExts {
		extSet[ext] = struct{}{}
	}
	err := filepath.Walk(dir.path, func(path string, file os.FileInfo, err error) error {
		if err != nil {
			s.logger.Warn(err)
		}
		if !file.IsDir() {
			ext := filepath.Ext(file.Name())
			if _, found := extSet[ext]; found {
				name := strings.TrimSuffix(file.Name(), ext)
				dir.Candidates = append(dir.Candidates, gameCandidate{name: name, extantion: ext, path: path})
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *scaner) rateCandidates(d *gameDir) {
	if len(d.Candidates) == 0 {
		return
	}
	applySizeScore(d)
	applyExeScore(d)
	applyNameScore(d, s.params.targetNames)
}

func (s *scaner) choseBestCandidate(d *gameDir) {
	if len(d.Candidates) == 0 {
		return
	}
	var bc gameCandidate
	for _, c := range d.Candidates {
		if c.matchScore > bc.matchScore {
			bc = c
		}
	}
	bc.name = namePrettifier(bc.name)
	d.bestCandidate = bc
}
