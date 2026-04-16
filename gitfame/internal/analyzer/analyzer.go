package analyzer

import (
	"path/filepath"

	"gitlab.com/slon/shad-go/gitfame/internal/git"
)

type Analyzer struct {
	repoPath     string
	revision     string
	useCommiter  bool
	excludeGlob  []string
	restrictGlob []string
	extensions   []string
}

type AnalyzeResult struct {
	Stats []PersonStat
}

type PersonStat struct {
	Name    string
	Lines   int
	Commits map[string]struct{}
	Files   map[string]struct{}
}

func NewAnalyzer(
	repoPath string,
	revision string,
	useCommiter bool,
	extensions []string,
	exclude []string,
	restrictTo []string) *Analyzer {
	return &Analyzer{
		repoPath:     repoPath,
		revision:     revision,
		useCommiter:  useCommiter,
		extensions:   extensions,
		excludeGlob:  exclude,
		restrictGlob: restrictTo,
	}
}

func (a *Analyzer) Analyze() (AnalyzeResult, error) {
	var authorStats = make(map[string]PersonStat)

	lsTree, err := git.LsTree(a.repoPath, a.revision, true)

	if err != nil {
		return AnalyzeResult{}, err
	}

	for _, lsTreeObj := range lsTree {
		if lsTreeObj.ObjectType != "blob" {
			continue
		}

		if a.isNotIncludedExtension(lsTreeObj.Path) {
			continue
		}

		if a.isInExcludeGlob(lsTreeObj.Path) {
			continue
		}

		if !a.isInRestrictToGlob(lsTreeObj.Path) {
			continue
		}

		groupsInfo, err := git.Blame(a.repoPath, lsTreeObj.Path, a.revision, a.useCommiter)

		if err != nil {
			return AnalyzeResult{}, err
		}

		if len(groupsInfo) == 0 {
			lastModification, err := git.FileLastLog(a.repoPath, lsTreeObj.Path, a.revision)
			if err != nil {
				return AnalyzeResult{}, err
			}
			updateAuthorStatsEmptyFile(lastModification, lsTreeObj.Path, authorStats, a.useCommiter)
		}

		updateAuthorStats(groupsInfo, lsTreeObj.Path, authorStats)
	}

	var authorStatsSlice []PersonStat

	for _, val := range authorStats {
		authorStatsSlice = append(authorStatsSlice, val)
	}

	return AnalyzeResult{authorStatsSlice}, nil
}

func updateAuthorStats(groupsInfo []git.GroupInfo, filepath string, authorStats map[string]PersonStat) {
	for _, val := range groupsInfo {
		stat, ok := authorStats[val.Author]

		if !ok {
			stat.Name = val.Author
			stat.Commits = make(map[string]struct{})
			stat.Files = make(map[string]struct{})
		}

		stat.Commits[val.Hash] = struct{}{}
		stat.Files[filepath] = struct{}{}
		stat.Lines += val.LinesNumber

		authorStats[val.Author] = stat
	}
}

func updateAuthorStatsEmptyFile(gitLogObject git.GitLogObject, filePath string, authorStats map[string]PersonStat, useCommiter bool) {
	author := gitLogObject.Author

	if useCommiter {
		author = gitLogObject.Commiter
	}

	stat, ok := authorStats[author]

	if !ok {
		stat.Name = author
		stat.Commits = make(map[string]struct{})
		stat.Files = make(map[string]struct{})
	}

	stat.Commits[gitLogObject.CommitHash] = struct{}{}
	stat.Files[filePath] = struct{}{}

	authorStats[author] = stat
}

func (a *Analyzer) isNotIncludedExtension(path string) bool {
	if a.extensions == nil {
		return false
	}

	for _, ext := range a.extensions {
		if filepath.Ext(path) == ext {
			return false
		}
	}

	return true
}

func (a *Analyzer) isInExcludeGlob(path string) bool {
	if a.excludeGlob == nil {
		return false
	}

	return isAnyMatchGlob(path, a.excludeGlob)
}

func (a *Analyzer) isInRestrictToGlob(path string) bool {
	if a.restrictGlob == nil {
		return true
	}

	return isAnyMatchGlob(path, a.restrictGlob)
}

func isAnyMatchGlob(path string, globs []string) bool {
	for _, glob := range globs {
		if ok, _ := filepath.Match(glob, path); ok {
			return true
		}
	}

	return false
}
