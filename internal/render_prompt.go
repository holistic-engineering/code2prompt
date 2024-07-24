package internal

import (
	"path/filepath"
	"strings"

	"github.com/aymerick/raymond"
)

const defaultTemplate = 
`# Project Structure
{{sourceTree}}

# Files
{{#each files}}
## {{this.Path}}
'''{{getFileExtension this.Path}}
{{this.Content}}
'''

{{/each}}
`

func RenderPrompt(files []FileInfo, sourceTree string, templateContent string) (string, error) {
	if templateContent == "" {
		templateContent = defaultTemplate
	}

	raymond.RegisterHelper("getFileExtension", func(path string) string {
		return strings.TrimPrefix(filepath.Ext(path), ".")
	})

	template, err := raymond.Parse(templateContent)
	if err != nil {
		return "", err
	}

	ctx := map[string]any{
		"files":      files,
		"sourceTree": sourceTree,
	}

	return template.Exec(ctx)
}