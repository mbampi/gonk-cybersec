package cronk

import (
	"log"
	"regexp"
	"strings"
)

// Routine represents a cron job with its description, time, and command.
type Routine struct {
	Description []string `json:"description"`
	Time        string   `json:"time"`
	Command     string   `json:"command"`
}

// Json represents the structure of the JSON output.
type Json struct {
	Intro    []string  `json:"intro"`
	Commands []Routine `json:"commands"`
	Outro    []string  `json:"outro"`
}

// CronToJson converts a cron file content into a structured JSON format.
func CronToJson(text string) Json {
	log.Println("Converting cron file to JSON")

	lines := strings.Split(text, "\n")
	commandIdx := getCommandLineIndices(lines)

	if len(commandIdx) == 0 {
		return Json{Intro: lines, Outro: []string{}}
	}

	intro, commandComments, outro := splitComments(lines, commandIdx)

	var routines []Routine
	for i, idx := range commandIdx {
		routines = append(routines, Routine{
			Description: commandComments[i],
			Time:        "", // Placeholder: real cron parsing would be added here
			Command:     strings.TrimSpace(lines[idx]),
		})
	}

	if intro == nil {
		intro = []string{}
	}
	if outro == nil {
		outro = []string{}
	}

	return Json{
		Intro:    intro,
		Commands: routines,
		Outro:    outro,
	}
}

// Helpers

func isCommand(line string) bool {
	matched, _ := regexp.MatchString(`^ *#`, line)
	return strings.TrimSpace(line) != "" && !matched
}

func getCommandLineIndices(lines []string) []int {
	var indices []int
	for i, line := range lines {
		if isCommand(line) {
			indices = append(indices, i)
		}
	}
	return indices
}

func splitComments(lines []string, commandIdx []int) ([]string, [][]string, []string) {
	endOfIntro := 0
	for i, line := range lines[:commandIdx[0]] {
		if strings.TrimSpace(line) == "" {
			endOfIntro = i
		}
	}

	intro := lines[:endOfIntro]
	if intro == nil {
		intro = []string{}
	}

	var commentBlocks [][]string
	starts := append([]int{endOfIntro}, commandIdx...)
	ends := append(commandIdx, len(lines))

	for i := range starts {
		start := starts[i] + 1
		end := ends[i]
		if start >= len(lines) || start >= end {
			commentBlocks = append(commentBlocks, []string{})
			continue
		}
		commentBlocks = append(commentBlocks, lines[start:end])
	}

	outro := commentBlocks[len(commentBlocks)-1]
	commentBlocks = commentBlocks[:len(commentBlocks)-1]

	// Remove trailing blank lines in outro
	var cleanOutro []string
	for _, line := range outro {
		if strings.TrimSpace(line) != "" {
			cleanOutro = append(cleanOutro, line)
		}
	}
	if cleanOutro == nil {
		cleanOutro = []string{}
	}

	return intro, commentBlocks, cleanOutro
}
