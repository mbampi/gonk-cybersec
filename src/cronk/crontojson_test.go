package cronk

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCronToJson(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Json
	}{
		{
			name: "No commands, only intro",
			input: `# only comments
# nothing else`,
			expected: Json{
				Intro:    []string{"# only comments", "# nothing else"},
				Commands: nil,
				Outro:    []string{},
			},
		},
		{
			name: "Intro, one command with comments",
			input: `# intro comment
# still intro

# command comment
0 0 * * * echo Hello World`,
			expected: Json{
				Intro: []string{"# intro comment", "# still intro"},
				Commands: []Routine{
					{
						Description: []string{"# command comment"},
						Time:        "",
						Command:     "0 0 * * * echo Hello World",
					},
				},
				Outro: []string{},
			},
		},
		{
			name:  "One command, no intro or outro",
			input: `0 0 * * * echo Only`,
			expected: Json{
				Intro: []string{},
				Commands: []Routine{
					{
						Description: []string{},
						Time:        "",
						Command:     "0 0 * * * echo Only",
					},
				},
				Outro: []string{},
			},
		},
		{
			name: "Multiple commands with outro",
			input: `# intro

# desc 1
0 * * * * echo One
# desc 2
30 * * * * echo Two

# outro`,
			expected: Json{
				Intro: []string{"# intro"},
				Commands: []Routine{
					{
						Description: []string{"# desc 1"},
						Time:        "",
						Command:     "0 * * * * echo One",
					},
					{
						Description: []string{"# desc 2"},
						Time:        "",
						Command:     "30 * * * * echo Two",
					},
				},
				Outro: []string{"# outro"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CronToJson(tc.input)
			if diff := cmp.Diff(tc.expected, got); diff != "" {
				t.Errorf("Mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}
