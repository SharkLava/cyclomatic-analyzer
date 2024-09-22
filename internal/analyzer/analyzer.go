package analyzer

import (
	"errors"
	"regexp"
)

// FileAnalysis represents analysis results for a single file.
type FileAnalysis struct {
	FilePath  string
	Functions []FunctionAnalysis
}

// FunctionAnalysis represents analysis results for a single function.
type FunctionAnalysis struct {
	Name           string
	Cyclomatic     int
	DecisionPoints []string
}

// Function represents a minimal version of a parsed function.
type Function struct {
	Name string
	Body string
}

// ErrUnmatchedBraces indicates unmatched braces during parsing.
var ErrUnmatchedBraces = errors.New("unmatched braces")

var decisionPointPatterns = []*regexp.Regexp{
	regexp.MustCompile(`\bif\b`),
	regexp.MustCompile(`\belse\s+if\b`),
	regexp.MustCompile(`\bfor\b`),
	regexp.MustCompile(`\bwhile\b`),
	regexp.MustCompile(`\bswitch\b`),
	regexp.MustCompile(`\bcase\b`),
	regexp.MustCompile(`\bcatch\b`),
	regexp.MustCompile(`&&`),
	regexp.MustCompile(`\|\|`),
	regexp.MustCompile(`\bdo\b`),
	regexp.MustCompile(`\?`),
}

// AnalyzeFunctions calculates cyclomatic complexity for a list of functions.
func AnalyzeFunctions(functions []Function) []FunctionAnalysis {
	analyses := make([]FunctionAnalysis, 0, len(functions))
	for _, fn := range functions {
		complexity, decisionPoints := calculateCyclomaticComplexity(fn.Body)
		analyses = append(analyses, FunctionAnalysis{
			Name:           fn.Name,
			Cyclomatic:     complexity,
			DecisionPoints: decisionPoints,
		})
	}
	return analyses
}

func calculateCyclomaticComplexity(code string) (int, []string) {
	var decisionPoints []string
	for _, regex := range decisionPointPatterns {
		matches := regex.FindAllString(code, -1)
		decisionPoints = append(decisionPoints, matches...)
	}
	return len(decisionPoints) + 1, decisionPoints
}
