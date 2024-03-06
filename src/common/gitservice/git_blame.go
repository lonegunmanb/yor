package gitservice

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/object"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/emirpasic/gods/queues/linkedlistqueue"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

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
	// Prepare the git command
	hash := c.Hash.String()
	cmd := exec.Command("git", "blame", hash, "--line-porcelain", "--", path)
	cmd.Dir = filepath.Dir(path)

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
	lines := linkedlistqueue.New()
	for _, line := range strings.Split(blameStr, "\n") {
		lines.Enqueue(line)
	}

	blameLines := make([]*git.Line, 0)
	var current *lineString
	for {
		line, ok := lines.Dequeue()
		if !ok {
			break
		}
		lineStr := line.(string)
		if strings.TrimSpace(lineStr) == "" {
			break
		}
		l := parseLine(lineStr, current)
		if l != current {
			if current != nil {
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
			}
			current = l
		}
	}

	return blameLines, nil
}

func parseLine(line string, current *lineString) *lineString {
	if strings.HasPrefix(line, "summary ") {

	} else if strings.HasPrefix(line, "author ") {
		current.AuthorString = line
	} else if strings.HasPrefix(line, "author-time ") {
		current.TimeString = line
	} else if strings.HasPrefix(line, "author-tz ") {
		current.TimezoneString = line
	} else if strings.HasPrefix(line, "committer ") {

	} else if strings.HasPrefix(line, "committer-mail ") {

	} else if strings.HasPrefix(line, "committer-time") {

	} else if strings.HasPrefix(line, "committer-tz ") {

	} else if strings.HasPrefix(line, "previous ") {

	} else {
		parts := strings.Split(line, " ")
		if len(parts) == 4 {
			newLs := &lineString{
				HashString: line,
			}
			return newLs
		}
		current.TextString = line
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
	prefix := "author "
	if !strings.HasPrefix(line, prefix) {
		return "", fmt.Errorf("invalid author line format")
	}
	return strings.TrimPrefix(line, prefix), nil
}

func parseHash(line string) (plumbing.Hash, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 4 {
		return plumbing.Hash{}, fmt.Errorf("invalid hash line format")
	}
	return plumbing.NewHash(parts[0]), nil
}
