package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/holistic-engineering/code2prompt/internal"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "code2prompt [path]",
	Short: "Generate LLM prompts from your codebase",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		exclude, _ := cmd.Flags().GetStringSlice("exclude")
		template, _ := cmd.Flags().GetString("template")
		outputFile, _ := cmd.Flags().GetString("output")
		jsonOutput, _ := cmd.Flags().GetBool("json")
		tokens, _ := cmd.Flags().GetBool("tokens")
		encoding, _ := cmd.Flags().GetString("encoding")
		maxFilesPerDir, _ := cmd.Flags().GetInt("max-files-per-dir")
		noSample, _ := cmd.Flags().GetBool("no-sample")

		files, err := internal.TraverseDirectory(path, exclude, maxFilesPerDir, noSample)
		if err != nil {
			fmt.Println("Error traversing directory:", err)
			os.Exit(1)
		}

		sourceTree := internal.GenerateSourceTree(files)
		prompt, err := internal.RenderPrompt(files, sourceTree.String(), template)
		if err != nil {
			fmt.Println("Error rendering prompt:", err)
			os.Exit(1)
		}

		var tokenCount int
		if tokens {
			tokenCount, err = internal.CountTokens(prompt, encoding)
			if err != nil {
				fmt.Println("Error counting tokens:", err)
			} else {
				fmt.Printf("Token count: %d\n", tokenCount)
			}
		}

		output := prompt
		if jsonOutput {
			outMap := map[string]any{
				"prompt": prompt,
			}

			if tokens {
				outMap["token_count"] = tokenCount
			}

			jsonData, _ := json.MarshalIndent(outMap, "", "  ")
			output = string(jsonData)
		}

		if outputFile != "" {
			err := os.WriteFile(outputFile, []byte(output), 0644)
			if err != nil {
				fmt.Println("Error writing to output file:", err)
				os.Exit(1)
			}
		} else {
			fmt.Println(output)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringSliceP("exclude", "e", []string{}, "Additional file patterns to exclude")
	rootCmd.Flags().StringP("template", "t", "", "Path to custom template file")
	rootCmd.Flags().StringP("output", "o", "", "Output file path")
	rootCmd.Flags().Bool("json", false, "Output as JSON")
	rootCmd.Flags().Bool("tokens", false, "Count tokens in resulting prompt")
	rootCmd.Flags().String("encoding", "cl100k_base", "Tokenizer encoding to use")
	rootCmd.Flags().IntP("max-files-per-dir", "m", 5, "Maximum number of files to include per directory")
	rootCmd.Flags().Bool("no-sample", false, "Include all files without sampling (overrides max-files-per-dir)")
}
