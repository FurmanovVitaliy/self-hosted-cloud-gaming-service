package scanner

type params struct {
	excludeDirs []string
	targetDirs  []string
	targetNames []string
	targetExts  []string
}

func CreateParams(excludeDirs, targetDirs, targetNames, targetExts []string) *params {
	return &params{
		excludeDirs: excludeDirs,
		targetDirs:  targetDirs,
		targetNames: targetNames,
		targetExts:  targetExts,
	}
}
