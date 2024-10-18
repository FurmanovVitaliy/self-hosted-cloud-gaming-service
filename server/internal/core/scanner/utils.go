package scanner

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

func namePrettifier(name string) string {
	wordsToRemove := []string{"dlc", "crack", "patch", "update", "fix", "repack", "multi", "rus", "eng", "engrus", "engrusmulti", "multie", "pack"}
	regexPattern := "\\b(" + strings.Join(wordsToRemove, "|") + ")\\b"

	name = regexp.MustCompile("([a-z])([A-Z])").ReplaceAllString(name, "$1 $2")
	name = regexp.MustCompile(`([a-zA-Z])\d{2}\b`).ReplaceAllString(name, "$1")

	prettyName := regexp.MustCompile(`\[[^\[\]]*\]`).ReplaceAllString(name, "")
	prettyName = regexp.MustCompile(`\([^()]*\)`).ReplaceAllString(prettyName, "")
	prettyName = regexp.MustCompile(`([^ ]) ([A-Z])`).ReplaceAllString(prettyName, "$1 $2")
	prettyName = regexp.MustCompile(`[^a-zA-Z0-9'"-: ]+`).ReplaceAllString(prettyName, "")
	prettyName = strings.ToLower(prettyName)
	prettyName = regexp.MustCompile(regexPattern).ReplaceAllString(prettyName, "")
	prettyName = regexp.MustCompile(`\s+`).ReplaceAllString(prettyName, " ")

	return prettyName
}

func notSystemDirectory(name string, excludeNames []string) bool {
	name = namePrettifier(name)
	bool := true
	for _, dir := range excludeNames {
		if strings.Contains(name, dir) {
			bool = false
		}
	}
	return bool
}

func levenshteinDistanceInPercent(s1, s2 string) int {
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

func platformByExtantion(extantion string) string {
	switch extantion {
	case ".exe":
		return "Windows"
	case ".nsp":
		return "Switch"
	case ".x86_64":
		return "Linux"
	default:
		return "unknown"
	}
}

func getFoldersHash(paths []string) (string, error) {
	var totalSize int64
	var lastModified time.Time

	// Сортируем пути
	sort.Strings(paths)

	for _, path := range paths {
		// Обходим все файлы в директории
		err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// Проверяем, что это файл, а не директория
			if !info.IsDir() {
				// Пропускаем файлы с расширением .ppbd
				if filepath.Ext(filePath) == ".ppbd" {
					return nil
				}
				if filepath.Ext(filePath) == ".ini" {
					return nil
				}
				// Увеличиваем общий размер файлов
				totalSize += info.Size()
				// Обновляем время последнего изменения, если оно новее
				if info.ModTime().After(lastModified) {
					lastModified = info.ModTime()
				}
			}
			return nil
		})

		if err != nil {
			return "", err
		}
	}

	// Формируем строку с данными о размере и времени последнего изменения
	data := fmt.Sprintf("%d-%s", totalSize, lastModified.String())
	// Создаем хеш MD5
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:]), nil
}
