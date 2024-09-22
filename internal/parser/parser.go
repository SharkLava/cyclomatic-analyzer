package parser

import (
	"github.com/SharkLava/cyclomatic-analyzer/internal/analyzer"
	"regexp"
)

// Function represents a C function with its name and body.
type Function struct {
	Name string
	Body string
}

var (
	funcRegex    = regexp.MustCompile(`(?m)^[a-zA-Z_][a-zA-Z0-9_]*\s+\**\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*\([^)]*\)\s*\{`)
	commentRegex = regexp.MustCompile(`("(?:\\.|[^"\\])*")|('(?:\\.|[^'\\])*')|(//.*?$)|(/\*.*?\*/)`)
)

// ParseFunctions parses the C code and extracts functions.
func ParseFunctions(code string) ([]Function, error) {
	code = removeCommentsAndStrings(code)
	matches := funcRegex.FindAllStringSubmatchIndex(code, -1)
	return extractFunctions(code, matches)
}

func extractFunctions(code string, matches [][]int) ([]Function, error) {
	var functions []Function

	for i, match := range matches {
		if len(match) < 4 {
			continue
		}
		nameStart, nameEnd := match[2], match[3]
		funcName := code[nameStart:nameEnd]

		bodyStart := match[1]
		body, err := extractBraces(code[bodyStart:])
		if err != nil {
			continue
		}

		functions = append(functions, Function{
			Name: funcName,
			Body: body,
		})

		if i+1 < len(matches) && matches[i+1][0] < bodyStart+len(body) {
			continue
		}
	}

	return functions, nil
}

func removeCommentsAndStrings(code string) string {
	return commentRegex.ReplaceAllString(code, "")
}

func extractBraces(code string) (string, error) {
	var braceCount int
	for i, char := range code {
		switch char {
		case '{':
			braceCount++
		case '}':
			braceCount--
			if braceCount == 0 {
				return code[:i+1], nil
			}
		}
	}
	return "", analyzer.ErrUnmatchedBraces
}
