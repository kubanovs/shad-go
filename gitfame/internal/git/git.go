package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type GroupInfo struct {
	Hash        string
	Author      string
	FilePath    string
	LinesNumber int
}

type LsTreeObject struct {
	AccessMode string
	ObjectType string
	Hash       string
	Path       string
}

type GitLogObject struct {
	Author     string
	Commiter   string
	CommitHash string
}

var ignoredPrefixes = []string{
	"author-mail ", "author-time ", "author-tz ", "committer-mail ", "committer-time ", "committer-tz ",
	"summary ",
}

func FileLastLog(repositoryPath, filePath, revision string) (GitLogObject, error) {
	args := []string{
		"log",
		revision,
		"-1",
		"--format=%an%n%cn%n%H",
		"--",
		filePath,
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = repositoryPath

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return GitLogObject{}, fmt.Errorf("error while executing %s: %s %v", strings.Join(cmd.Args, " "), &stderr, err)
	}

	parts := strings.Split(strings.TrimRight(stdout.String(), "\n"), "\n")
	if len(parts) < 3 {
		return GitLogObject{}, fmt.Errorf("unexpected git log output")
	}

	return GitLogObject{Author: parts[0], Commiter: parts[1], CommitHash: parts[2]}, nil
}

func Blame(repositoryPath, filepath, revision string, useCommiter bool) ([]GroupInfo, error) {
	args := []string{
		"blame",
		"--porcelain",
		revision,
		filepath,
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = repositoryPath

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error while executing %s: %s %v", strings.Join(cmd.Args, " "), &stderr, err)
	}

	return parseBlameOutput(stdout.Bytes(), useCommiter)
}

func LsTree(repositoryPath, revision string, useRecursion bool) ([]LsTreeObject, error) {
	args := []string{
		"-C",
		repositoryPath,
		"ls-tree",
		revision,
		"--full-tree",
	}

	if useRecursion {
		args = append(args, "-r")
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = repositoryPath

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error while run %s: %s %v", strings.Join(cmd.Args, " "), &stderr, err)
	}

	return parceLsTreeOutput(&stdout), nil
}

func parceLsTreeOutput(stdout *bytes.Buffer) []LsTreeObject {
	var paths []LsTreeObject

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "\t")

		part1 := parts[0]
		path := parts[1]

		parts2 := strings.Split(part1, " ")

		accessMode := parts2[0]
		objectType := parts2[1]
		hash := parts2[2]

		paths = append(paths, LsTreeObject{AccessMode: accessMode, ObjectType: objectType, Hash: hash, Path: path})
	}

	return paths
}

func parseBlameOutput(data []byte, useCommiter bool) ([]GroupInfo, error) {
	groups := make(map[string]GroupInfo)
	scanner := bufio.NewScanner(bytes.NewReader(data))

	var currentHash string

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "\t"):
			continue

		case isIgnoredLine(line):
			continue

		case strings.HasPrefix(line, "author "):
			if !useCommiter {
				group := groups[currentHash]
				group.Author = strings.TrimPrefix(line, "author ")
				groups[currentHash] = group
			}
		case strings.HasPrefix(line, "committer "):
			if useCommiter {
				group := groups[currentHash]
				group.Author = strings.TrimPrefix(line, "committer ")
				groups[currentHash] = group
			}
		case strings.HasPrefix(line, "filename "):
			group := groups[currentHash]
			group.FilePath = strings.TrimPrefix(line, "filename ")
			groups[currentHash] = group

		default:
			hash, linesCount, err := parseCommitLine(line)
			if err != nil {
				continue // не commit-строка
			}

			currentHash = hash
			group := groups[currentHash]
			group.Hash = hash
			group.LinesNumber += linesCount
			groups[currentHash] = group
		}
	}

	result := make([]GroupInfo, 0, len(groups))
	for _, group := range groups {
		result = append(result, group)
	}
	return result, nil
}

func parseCommitLine(line string) (hash string, linesCount int, err error) {
	parts := strings.Split(line, " ")
	if len(parts) != 4 {
		return "", 0, fmt.Errorf("not a commit line")
	}

	linesCount, err = strconv.Atoi(parts[3])
	if err != nil {
		return "", 0, fmt.Errorf("parse lines count %q: %w", parts[3], err)
	}

	return parts[0], linesCount, nil
}

func isIgnoredLine(line string) bool {
	for _, prefix := range ignoredPrefixes {
		if strings.HasPrefix(line, prefix) {
			return true
		}
	}
	return false
}
