package gitservice

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func timeParse(layout string, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

func TestParseBlameResult(t *testing.T) {
	cases := []struct {
		desc     string
		input    string
		expected []*git.Line
	}{
		{
			desc:     "no line",
			input:    ``,
			expected: []*git.Line{},
		},
		{
			desc: "single line",
			input: fmt.Sprintf(`cc26c4d06f8e57f8da2a8b7d008142b5026556cb 1 1 1
author hezijie
author-mail <lonegunmanb@hotmail.com>
author-time 1709707964
author-tz +0800
committer hezijie
committer-mail <lonegunmanb@hotmail.com>
committer-time 1709707964
committer-tz +0800
summary multiple lines test
previous fb8e4e2e0b8a874c18ad32a319a8d79c37519d2a main.tf
filename main.tf
%sdata "azurerm_resource_group" "main"  {`, "\t"),
			expected: []*git.Line{
				{
					Author: "hezijie",
					Text:   `        data "azurerm_resource_group" "main"  {`,
					//"Mon, 02 Jan 2006 15:04:05 -0700"
					Date: timeParse(time.RFC1123Z, "Wed, 06 Mar 2024 14:52:44 +0800").UTC(),
					Hash: plumbing.NewHash("cc26c4d06f8e57f8da2a8b7d008142b5026556cb"),
				},
			},
		},
		{
			desc: "newline tail",
			input: fmt.Sprintf(`cc26c4d06f8e57f8da2a8b7d008142b5026556cb 1 1 1
author hezijie
author-mail <lonegunmanb@hotmail.com>
author-time 1709707964
author-tz +0800
committer hezijie
committer-mail <lonegunmanb@hotmail.com>
committer-time 1709707964
committer-tz +0800
summary multiple lines test
previous fb8e4e2e0b8a874c18ad32a319a8d79c37519d2a main.tf
filename main.tf
%sdata "azurerm_resource_group" "main"  {
`, "\t"),
			expected: []*git.Line{
				{
					Author: "hezijie",
					Text:   `        data "azurerm_resource_group" "main"  {`,
					//"Mon, 02 Jan 2006 15:04:05 -0700"
					Date: timeParse(time.RFC1123Z, "Wed, 06 Mar 2024 14:52:44 +0800").UTC(),
					Hash: plumbing.NewHash("cc26c4d06f8e57f8da2a8b7d008142b5026556cb"),
				},
			},
		},
		{
			desc: "multiple line",
			input: fmt.Sprintf(`cc26c4d06f8e57f8da2a8b7d008142b5026556cb 1 1 1
author hezijie
author-mail <lonegunmanb@hotmail.com>
author-time 1709707964
author-tz +0800
committer hezijie
committer-mail <lonegunmanb@hotmail.com>
committer-time 1709707964
committer-tz +0800
summary multiple lines test
previous fb8e4e2e0b8a874c18ad32a319a8d79c37519d2a main.tf
filename main.tf
%sdata "azurerm_resource_group" "main"  {
e6b0bff7323580793611d8f68db93deddd6f2a46 2 2 1
author Yuping Wei
author-mail <56525716+yupwei68@users.noreply.github.com>
author-time 1581072126
author-tz +0800
committer GitHub
committer-mail <noreply@github.com>
committer-time 1581072126
committer-tz +0800
summary terraform version upgrade and code reorg (#39)
previous d365197adaefbd87f269772556793f2e46f00c41 main.tf
filename main.tf
%sname = var.resource_group_name`, "\t", "\t"),
			expected: []*git.Line{
				{
					Author: "hezijie",
					Text:   `        data "azurerm_resource_group" "main"  {`,
					Date:   timeParse(time.RFC1123Z, "Wed, 06 Mar 2024 14:52:44 +0800").UTC(),
					Hash:   plumbing.NewHash("cc26c4d06f8e57f8da2a8b7d008142b5026556cb"),
				},
				{
					Author: "Yuping Wei",
					Text:   `          name = var.resource_group_name`,
					Date:   timeParse(time.RFC1123Z, "Fri, 07 Feb 2020 18:42:06 +0800").UTC(),
					Hash:   plumbing.NewHash("e6b0bff7323580793611d8f68db93deddd6f2a46"),
				},
			},
		},
	}

	for _, c := range cases {
		cc := c
		t.Run(cc.desc, func(t *testing.T) {
			g := &GitService{}
			result, err := g.ParseBlameResult(cc.input)
			if err != nil {
				t.Fatalf("ParseBlameResult returned an error: %v", err)
			}

			if len(result) != len(cc.expected) {
				t.Fatalf("Expected %d lines, got %d", len(cc.expected), len(result))
			}

			for i, line := range result {
				assert.Equal(t, *cc.expected[i], *line)
			}
		})
	}

}
