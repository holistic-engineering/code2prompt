package internal

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"math/rand"

	"github.com/gobwas/glob"
)

type FileInfo struct {
	Path    string
	Content string
}

const MaxFilesPerDirectory = 5

func TraverseDirectory(root string, additionalExcludes []string) ([]FileInfo, error) {
	var files []FileInfo
	
	gitignorePatterns, err := readGitignore(root)
	if err != nil {
		return nil, err
	}

	defaultExcludes := []string{
		"node_modules/**",
		"vendor/**",
		"build/**",
		"dist/**",
		"*.exe",
		"*.dll",
		"*.so",
		"*.dylib",
		"*.class",
		"*.jar",
		"*.war",
		"*.ear",
		"*.zip",
		"*.tar.gz",
		"*.rar",
		"*.log",
		".git/**",
		".svn/**",
		".hg/**",
		"*.pdf",
		"*.png",
		"*.jpg",
		"*.jpeg",
		"*.gif",
		"*.bmp",
		"*.tiff",
		"*.ico",
		"*.svg",
		"*.webp",
		"*.min.js",
		"*.min.css",
		"*.lock",
		"*test*/**",        // Exclude directories with 'test' in the name
		"**/*test*.{js,py,go,java,cs,ts,cpp,c,rb}", // Exclude test files with common extensions
		"**/*spec*.{js,ts}", // Exclude spec files (common in JavaScript/TypeScript projects)
		"**/__tests__/**",   // Exclude __tests__ directories (common in React projects)
	}

	excludePatterns := append(gitignorePatterns, defaultExcludes...)
	excludePatterns = append(excludePatterns, additionalExcludes...)

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
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

	return SampleFiles(files), nil
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
			patterns = append(patterns, pattern)
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

func SampleFiles(files []FileInfo) []FileInfo {
	filesByDir := make(map[string][]FileInfo)
	
	for _, file := range files {
		dir := filepath.Dir(file.Path)
		filesByDir[dir] = append(filesByDir[dir], file)
	}

	var sampledFiles []FileInfo
	for _, dirFiles := range filesByDir {
		if len(dirFiles) > MaxFilesPerDirectory {
			rand.Shuffle(len(dirFiles), func(i, j int) {
				dirFiles[i], dirFiles[j] = dirFiles[j], dirFiles[i]
			})
			dirFiles = dirFiles[:MaxFilesPerDirectory]
		}
		sampledFiles = append(sampledFiles, dirFiles...)
	}

	return sampledFiles
}