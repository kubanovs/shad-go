package gfcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	language "gitlab.com/slon/shad-go/gitfame/configs"
	"gitlab.com/slon/shad-go/gitfame/internal/analyzer"
	"gitlab.com/slon/shad-go/gitfame/internal/formatter"
)

var (
	revision    string
	repoPath    string
	format      string
	orderBy     string
	useCommiter bool
	extensions  []string
	exclude     []string
	restrictTo  []string
	languages   []string
)

var columns = map[string]struct{}{
	"lines":   struct{}{},
	"commits": struct{}{},
	"files":   struct{}{},
}

func Main() int {
	rootCmd := &cobra.Command{
		Use:  "gitfame",
		RunE: process,
	}

	rootCmd.Flags().StringVar(&revision, "revision", "HEAD", "revision of git repository")
	rootCmd.Flags().StringVar(&repoPath, "repository", ".", "path to git repository")
	rootCmd.Flags().StringVar(&format, "format", "tabular", "output format")
	rootCmd.Flags().StringVar(&orderBy, "order-by", "lines", "order-by column name")
	rootCmd.Flags().BoolVar(&useCommiter, "use-committer", false, "use commiter as author")
	rootCmd.Flags().StringSliceVar(&extensions, "extensions", nil, "included extensions (ex. .go)")
	rootCmd.Flags().StringSliceVar(&exclude, "exclude", nil, "exclude glob patterns")
	rootCmd.Flags().StringSliceVar(&restrictTo, "restrict-to", nil, "restrict to glob patterns")
	rootCmd.Flags().StringSliceVar(&languages, "languages", nil, "allowed languages")
	err := rootCmd.Execute()

	if err != nil {
		return 1
	}
	return 0
}

func process(cmd *cobra.Command, args []string) error {

	if _, ok := columns[orderBy]; !ok {
		return fmt.Errorf("not valid orderBy arg: %s", orderBy)
	}

	var resultExtensions []string = nil

	if extensions != nil || languages != nil {
		resultExtensions = []string{}

		resultExtensions = append(resultExtensions, extensions...)
		resultExtensions = append(resultExtensions, language.New().GetLanguagesExtensions(languages)...)
	}

	analyzer := analyzer.NewAnalyzer(repoPath, revision, useCommiter, resultExtensions, exclude, restrictTo)

	analyzeResult, err := analyzer.Analyze()

	if err != nil {
		return err
	}

	statStr, err := formatter.Format(analyzeResult.Stats, orderBy, format)

	if err != nil {
		return err
	}

	fmt.Print(statStr)

	return nil
}
