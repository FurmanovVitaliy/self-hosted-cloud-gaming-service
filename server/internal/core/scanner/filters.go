package scanner

import (
	"os"

	"github.com/sahilm/fuzzy"
)

func applySizeScore(d *gameDir) {
	var maxFileSize int64
	for i, c := range d.Candidates {
		fileInfo, err := os.Stat(c.path)
		if err != nil {
			continue
		}
		if fileInfo.Size() > maxFileSize {
			maxFileSize = fileInfo.Size()
			d.Candidates[i].matchScore += 4
		}
	}
}

func applyExeScore(d *gameDir) {
	for i, c := range d.Candidates {
		if c.extantion == ".exe" {
			if _, err := os.Stat(c.path + ".ppdb"); err == nil {
				d.Candidates[i].matchScore += 4
			}
		}
	}
}

func applyNameScore(d *gameDir, targetNames []string) {
	for i, c := range d.Candidates {

		if c.name == "launcher" {
			d.Candidates[i].matchScore += 22
			d.Candidates[i].name = d.name
		}

		c.name = namePrettifier(c.name)
		matches := fuzzy.Find(c.name, targetNames)

		if len(matches) > 0 {
			d.Candidates[i].name = targetNames[matches[0].Index]
			d.Candidates[i].matchScore += 8
		}

		levMatch := levenshteinDistanceInPercent(c.name, namePrettifier(d.name))
		switch {
		case levMatch > 80:
			d.Candidates[i].matchScore += 8
		case levMatch > 60:
			d.Candidates[i].matchScore += 6
		case levMatch > 40:
			d.Candidates[i].matchScore += 4
		case levMatch > 20:
			d.Candidates[i].matchScore += 2
		}
	}
}
