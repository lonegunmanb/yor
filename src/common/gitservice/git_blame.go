package gitservice

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/object"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var sha1Regex = regexp.MustCompile("^[a-fA-F0-9]{40}")

type lineString struct {
	// Author is the email address of the last author that modified the line.
	AuthorString string
	// Text is the original text of the line.
	TextString string
	// Date is when the original text of the line was introduced
	TimeString     string
	TimezoneString string
	// Hash is the commit hash that introduced the original line
	HashString string
}

func (g *GitService) Blame(c *object.Commit, path string) (*git.BlameResult, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) && !filepath.IsAbs(path) {
		path = filepath.Join(g.gitRootDir, path)
	}
	path = filepath.Clean(path)
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	// Prepare the git command
	hash := c.Hash.String()
	cmd := exec.Command("git", "blame", hash, "--line-porcelain", "--", filepath.Base(path))
	cmd.Dir = filepath.Dir(absPath)

	// Run the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute git blame: %w", err)
	}

	// Convert the output to a string
	blameStr := string(output)

	// Parse the blame result
	blameLines, err := g.ParseBlameResult(blameStr)
	if err != nil {
		return nil, err
	}

	// Create the blame result
	blameResult := &git.BlameResult{
		Path:  path,
		Rev:   c.Hash,
		Lines: blameLines,
	}

	return blameResult, nil
}

func (g *GitService) ParseBlameResult(blameStr string) ([]*git.Line, error) {
	linesSlice := strings.Split(blameStr, "\n")

	blameLines := make([]*git.Line, 0)
	current := new(lineString)
	for _, lineStr := range linesSlice {
		l := parseLine(lineStr, current)
		if l != current {
			author, err := parseAuthor(current.AuthorString)
			if err != nil {
				return nil, err
			}
			authorTime, err := parseTime(current.TimeString, current.TimezoneString)
			if err != nil {
				return nil, err
			}
			lineHash, err := parseHash(current.HashString)
			if err != nil {
				return nil, err
			}
			text := current.TextString
			blameLines = append(blameLines, &git.Line{
				Author: author,
				Text:   text,
				Date:   authorTime,
				Hash:   lineHash,
			})
			current = l
		}
	}

	return blameLines, nil
}

func parseLine(line string, current *lineString) *lineString {
	if strings.HasPrefix(line, "author-mail") {
		current.AuthorString = line
	} else if strings.HasPrefix(line, "author-time ") {
		current.TimeString = line
	} else if strings.HasPrefix(line, "author-tz ") {
		current.TimezoneString = line
	} else if strings.HasPrefix(line, "\t") {
		current.TextString = strings.TrimPrefix(line, "\t")
		return new(lineString)
	} else {
		parts := strings.Split(line, " ")
		if len(parts) > 1 && sha1Regex.MatchString(parts[0]) {
			current.HashString = line
		}
	}
	return current
}

func parseTime(authorTimeLine string, authorTimezoneLine string) (time.Time, error) {
	authorTimePrefix := "author-time "
	if !strings.HasPrefix(authorTimeLine, authorTimePrefix) {
		return time.Time{}, fmt.Errorf("invalid author time format")
	}
	authorTime := strings.TrimPrefix(authorTimeLine, authorTimePrefix)

	authorTimezonePrefix := "author-tz "
	if !strings.HasPrefix(authorTimezoneLine, authorTimezonePrefix) {
		return time.Time{}, fmt.Errorf("invalid author timezone format")
	}
	authorTimezone := strings.TrimPrefix(authorTimezoneLine, authorTimezonePrefix)
	// Parse the Unix timestamp
	unixTime, err := strconv.ParseInt(authorTime, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid author time format")
	}

	// Convert the Unix timestamp to a time.Time object
	t := time.Unix(unixTime, 0)

	// Parse the timezone string
	if len(authorTimezone) != 5 {
		return time.Time{}, fmt.Errorf("invalid author timezone format")
	}

	sign := authorTimezone[:1]
	hours, err := strconv.Atoi(authorTimezone[1:3])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid author timezone format")
	}

	minutes, err := strconv.Atoi(authorTimezone[3:])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid author timezone format")
	}

	// Calculate the timezone offset in minutes
	offset := hours*60 + minutes
	if sign == "-" {
		offset = -offset
	}

	// Convert the time.Time object to the timezone
	t = t.In(time.UTC)

	return t, nil
}

func parseAuthor(line string) (string, error) {
	prefix := "author-mail "
	if !strings.HasPrefix(line, prefix) {
		return "", fmt.Errorf("invalid author line format")
	}
	line = strings.TrimPrefix(line, prefix)
	line = strings.TrimPrefix(line, "<")
	line = strings.TrimSuffix(line, ">")
	return line, nil
}

func parseHash(line string) (plumbing.Hash, error) {
	parts := strings.Split(line, " ")
	if len(parts) < 1 {
		return plumbing.Hash{}, fmt.Errorf("invalid hash line format")
	}
	return plumbing.NewHash(parts[0]), nil
}
