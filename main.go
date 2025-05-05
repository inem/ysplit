package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"strings"
)

func main() {
	asJSON := flag.Bool("json", false, "output JSON instead of YAML")
	flag.Parse()

	raw, _ := io.ReadAll(os.Stdin)

	parentKey, blocks := splitIntoBlocks(raw)
	if parentKey == "" {
		exit(fmt.Errorf("no top-level key found"))
	}

	// parse each block
	var list []any
	for _, b := range blocks {
		var m any
		if err := yaml.NewDecoder(bytes.NewReader([]byte(b))).Decode(&m); err != nil {
			exit(err)
		}
		list = append(list, m)
	}

	out := map[string]any{parentKey: list}

	if *asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(out)
	} else {
		yaml.NewEncoder(os.Stdout).Encode(out)
	}
}

// ----- helpers -------------------------------------------------------------

func splitIntoBlocks(src []byte) (parent string, blocks []string) {
	// First, split into lines to find the parent key
	lines := bytes.Split(src, []byte("\n"))
	if len(lines) == 0 {
		return "", nil
	}

	// Find the parent key (first line without indentation)
	parentLine := ""
	var contentLines []string
	for _, line := range lines {
		l := string(line)
		if len(l) > 0 && !isWhitespace(l) && parentLine == "" {
			parentLine = strings.TrimSuffix(strings.TrimSpace(l), ":")
			continue
		}
		if parentLine != "" {
			contentLines = append(contentLines, l)
		}
	}

	if parentLine == "" {
		return "", nil
	}

	// Split into blocks by empty lines
	var result []string
	var currentBlock []string
	expectNewBlock := true

	for _, line := range contentLines {
		if strings.TrimSpace(line) == "" {
			if len(currentBlock) > 0 {
				// If the block is not empty - add it to the result and start a new one
				yamlBlock := processBlock(currentBlock)
				if yamlBlock != "" {
					result = append(result, yamlBlock)
				}
				currentBlock = nil
			}
			expectNewBlock = true
			continue
		}

		if expectNewBlock && strings.TrimSpace(line) != "" {
			// Start a new block
			currentBlock = []string{line}
			expectNewBlock = false
		} else {
			// Add line to the current block
			currentBlock = append(currentBlock, line)
		}
	}

	// Don't forget the last block
	if len(currentBlock) > 0 {
		yamlBlock := processBlock(currentBlock)
		if yamlBlock != "" {
			result = append(result, yamlBlock)
		}
	}

	return parentLine, result
}

// Process a block of lines, preserving the correct nesting structure
func processBlock(lines []string) string {
	if len(lines) == 0 {
		return ""
	}

	// Find the base indentation (minimum for non-empty lines)
	baseIndent := -1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		ind := getIndent(line)
		if baseIndent == -1 || ind < baseIndent {
			baseIndent = ind
		}
	}

	if baseIndent <= 0 {
		// If there's no indentation, this is an invalid block
		return ""
	}

	// Create a new block with proper indentation
	var resultLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			resultLines = append(resultLines, "")
		} else {
			// Preserve relative indentation
			ind := getIndent(line)
			newIndent := ind - baseIndent

			// Format the line with the required indentation
			trimmedLine := strings.TrimLeft(line, " \t")
			resultLines = append(resultLines, strings.Repeat(" ", newIndent) + trimmedLine)
		}
	}

	return strings.Join(resultLines, "\n")
}

// Get the number of spaces at the beginning of a string
func getIndent(s string) int {
	return len(s) - len(strings.TrimLeft(s, " \t"))
}

// Check if a string consists only of whitespace characters
func isWhitespace(s string) bool {
	return strings.TrimSpace(s) == ""
}

func trimParentIndent(s string) string {
	return strings.TrimLeft(s, " ")
}

func exit(err error) {
	fmt.Fprintln(os.Stderr, "ysplit:", err)
	os.Exit(1)
}
