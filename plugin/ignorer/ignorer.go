package ignorer

import (
	"bufio"
	"io"

	gitignore "github.com/sabhiram/go-gitignore"
)

// Ignorer wraps an io.Reader containing ignore patterns.
type Ignorer interface {
	ShouldIgnore(changedFiles ...string) bool
}

type ignorer struct {

	// GitIgnorer contains the ignore patterns.
	GitIgnorer *gitignore.GitIgnore
}

// New returns a instance of DroneIgnore
func New(ignoreExpr io.Reader) (Ignorer, error) {

	var lines []string
	scanner := bufio.NewScanner(ignoreExpr)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	gitignore, err := gitignore.CompileIgnoreLines(lines...)

	if err != nil {
		return nil, err
	}

	return &ignorer{gitignore}, nil
}

// ShouldIgnore takes a list of changed files and returns 'true' if it should ignore
// these files, i.e., if the ignore expression matches all files in the list.
func (i *ignorer) ShouldIgnore(changedFiles ...string) bool {

	// If we DO NOT match any path, return false
	for _, file := range changedFiles {
		if !i.GitIgnorer.MatchesPath(file) {
			return false
		}
	}

	// If all patches are matched, we can return true
	return true
}
