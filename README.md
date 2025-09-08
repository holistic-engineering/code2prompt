# Code2Prompt

Code2Prompt is a powerful command-line tool that converts your entire codebase into a comprehensive prompt for Large Language Models (LLMs). It's designed to provide extensive context about your project, making it easier for LLMs to understand and assist with your codebase.

## Features

- **Codebase Traversal**: Automatically walks through your project directory, respecting .gitignore rules.
- **Source Tree Generation**: Creates a hierarchical representation of your project structure.
- **Smart File Sampling**: Configurable file sampling per directory to manage large codebases efficiently.
- **Prompt Templating**: Customizable prompt generation using Handlebars templates.
- **Token Counting**: Calculates the token count of the generated prompt for compatibility with various LLM token limits.
- **Multiple Output Formats**: Supports plain text and JSON output.
- **Flexible File Filtering**: Exclude additional files using glob patterns.

## Installation

### Prerequisites

- Go 1.22 or later

### Building from Source

1. Clone the repository:
   ```
   git clone https://github.com/holistic-engineering/code2prompt.git
   cd code2prompt
   ```

2. Build the project:
   ```
   make build
   ```

This will create a `code2prompt` binary in your current directory.

## Usage

Basic usage:

```
code2prompt [flags] <path>
```

### Flags

- `-e, --exclude`: File patterns to exclude (comma-separated)
- `-t, --template`: Path to custom template file
- `-o, --output`: Output file path
- `--tokens`: Count tokens in the generated prompt
- `--json`: Output as JSON
- `--encoding`: Tokenizer encoding to use (default "cl100k_base")
- `-m, --max-files-per-dir`: Maximum number of files to include per directory (default 5)
- `--no-sample`: Include all files without sampling (overrides max-files-per-dir)

### Examples

1. Generate a prompt for a project:
   ```
   code2prompt /path/to/your/project
   ```

2. Generate a prompt, excluding certain file types:
   ```
   code2prompt --exclude "*.md,*.txt" /path/to/your/project
   ```

3. Use a custom template and save the output:
   ```
   code2prompt -t custom_template.hbs -o output.md /path/to/your/project
   ```

4. Generate JSON output with token count:
   ```
   code2prompt --json --tokens /path/to/your/project
   ```

5. Include all files without sampling:
   ```
   code2prompt --no-sample /path/to/your/project
   ```

6. Set a higher file limit per directory:
   ```
   code2prompt --max-files-per-dir 10 /path/to/your/project
   ```

7. Combine multiple options:
   ```
   code2prompt --no-sample --exclude "*.log,*.tmp" --json /path/to/your/project
   ```

### File Sampling

By default, Code2Prompt limits the number of files included per directory to 5 to keep the output manageable for large codebases. Files are selected alphabetically (deterministic) when the limit is exceeded.

**Key behaviors:**
- **Default limit**: 5 files per directory
- **Selection method**: Alphabetical sorting (ensures consistent results)
- **Customization**: Use `--max-files-per-dir` to change the limit
- **Disable sampling**: Use `--no-sample` to include all files

**When to use each option:**
- **Default behavior**: Good for exploring large codebases
- **`--max-files-per-dir N`**: When you need more context but still want limits
- **`--no-sample`**: When you need complete codebase coverage (recommended for smaller projects or when specific files are being excluded)

### Exclude Patterns

The `--exclude` flag supports glob patterns for flexible file filtering. **Important**: Always quote patterns containing wildcards to prevent shell expansion.

**Pattern Examples:**
```bash
# ✅ Correct: Quoted patterns
code2prompt --exclude '*.log' /path/to/project
code2prompt --exclude '.godot/editor/*' /path/to/project
code2prompt --exclude '.godot/editor/' /path/to/project

# ❌ Incorrect: Unquoted wildcards (causes shell expansion errors)
code2prompt --exclude *.log /path/to/project  # Error: shell expands before tool sees it
code2prompt --exclude .godot/editor/* /path/to/project  # Error: "accepts 1 arg(s), received N"
```

**Pattern Types:**
- **File extensions**: `'*.log'`, `'*.tmp'`, `'*.bin'`
- **Directory contents**: `'.godot/editor/*'` (all files in directory)
- **Entire directories**: `'.godot/editor/'` (directory and all contents)
- **Nested patterns**: `'**/.godot/**'` (all .godot directories anywhere)

**Multiple patterns:**
```bash
code2prompt --exclude '*.log' --exclude '.godot/editor/*' --exclude 'build/' /path/to/project
```

### Custom Templates

You can create custom templates for the prompt using the [HandlebarsJS](https://handlebarsjs.com/) templating engine. The following list shows all available variables and helper functions. 

- `sourceTree`: tree representation of the project structure
- `files`: list of all filtered project files 
- `this.Path`: the individual file's path
- `this.Content`: the individual file's content
- `getFileExtension`: gets the file extension for a given a file path

#### Default Temlate

```Markdown
# Project Structure
{{sourceTree}}

# Files
{{#each files}}
## {{this.Path}}
'''{{getFileExtension this.Path}}
{{this.Content}}
'''

{{/each}}
```

## Development

### Prerequisites

- Go 1.22 or later
- Make

### Setup

1. Clone the repository:
   ```
   git clone https://github.com/holistic-engineering/code2prompt.git
   cd code2prompt
   ```

2. Install dependencies:
   ```
   make deps
   ```

### Building

To build the project:

```
make build
```

### Clean Up

To clean up built binaries:

```
make clean
```

## Contributing

Contributions to Code2Prompt are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) for CLI interface
- [Viper](https://github.com/spf13/viper) for configuration management
- [Raymond](https://github.com/aymerick/raymond) for Handlebars templating
- [TikToken](https://github.com/tiktoken-go/tokenizer) for token counting