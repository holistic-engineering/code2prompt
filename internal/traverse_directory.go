package internal

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gobwas/glob"
)

type FileInfo struct {
	Path    string
	Content string
}

// Removed MaxFilesPerDirectory constant - now configurable via command line flags

func TraverseDirectory(root string, additionalExcludes []string, maxFilesPerDir int, noSample bool) ([]FileInfo, error) {
	var files []FileInfo

	gitignorePatterns, err := readGitignore(root)
	if err != nil {
		return nil, err
	}

	defaultExcludes := []string{
		"**node_modules/**",
		"**vendor/**",
		"**build/**",
		"**dist/**",
		"**.exe",
		"**.dll",
		"**.so",
		"**.dylib",
		"**.class",
		"**.jar",
		"**.war",
		"**.ear",
		"**.zip",
		"**.tar.gz",
		"**.rar",
		"**.log",
		"**.git/**",
		"**.svn/**",
		"**.hg/**",
		"**.pdf",
		"**.png",
		"**.jpg",
		"**.jpeg",
		"**.gif",
		"**.bmp",
		"**.tiff",
		"**.ico",
		"**.svg",
		"**.webp",
		"**.min.js",
		"**.min.css",
		"**.lock",
		"**test*/**",
		"**test*.{js,py,go,java,cs,ts,cpp,c,rb}",
		"**spec*.{js,ts}",
		"**__tests__/**",
	}

	excludePatterns := append(gitignorePatterns, defaultExcludes...)
	excludePatterns = append(excludePatterns, additionalExcludes...)

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			if shouldExclude(relPath, excludePatterns) {
				return filepath.SkipDir
			}
			return nil
		}

		if !shouldExclude(relPath, excludePatterns) {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			files = append(files, FileInfo{
				Path:    relPath,
				Content: string(content),
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if noSample {
		return files, nil
	}
	return SampleFiles(files, maxFilesPerDir), nil
}

func readGitignore(root string) ([]string, error) {
	gitignorePath := filepath.Join(root, ".gitignore")
	patterns := []string{}

	file, err := os.Open(gitignorePath)
	if os.IsNotExist(err) {
		return patterns, nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pattern := strings.TrimSpace(scanner.Text())
		if pattern != "" && !strings.HasPrefix(pattern, "#") {
			patterns = append(patterns, "**"+pattern)
		}
	}

	return patterns, scanner.Err()
}

func shouldExclude(path string, patterns []string) bool {
	for _, pattern := range patterns {
		g, err := glob.Compile(pattern)
		if err != nil {
			continue
		}
		if g.Match(path) {
			return true
		}
	}
	return false
}

func SampleFiles(files []FileInfo, maxFilesPerDir int) []FileInfo {
	filesByDir := make(map[string][]FileInfo)

	for _, file := range files {
		dir := filepath.Dir(file.Path)
		filesByDir[dir] = append(filesByDir[dir], file)
	}

	var sampledFiles []FileInfo
	for _, dirFiles := range filesByDir {
		if len(dirFiles) > maxFilesPerDir {
			// Sort files to ensure deterministic behavior instead of random sampling
			sort.Slice(dirFiles, func(i, j int) bool {
				return dirFiles[i].Path < dirFiles[j].Path
			})
			dirFiles = dirFiles[:maxFilesPerDir]
		}
		sampledFiles = append(sampledFiles, dirFiles...)
	}

	return sampledFiles
}
