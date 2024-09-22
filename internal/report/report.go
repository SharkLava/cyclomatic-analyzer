package report

import (
	"encoding/json"
	"github.com/SharkLava/cyclomatic-analyzer/internal/analyzer"
)

// Report represents the complete analysis report.
type Report struct {
	TotalFiles      int          `json:"total_files"`
	TotalFunctions  int          `json:"total_functions"`
	TotalComplexity int          `json:"total_complexity"`
	Files           []FileReport `json:"files"`
}

// FileReport represents the analysis of a single file.
type FileReport struct {
	FilePath       string                      `json:"file_path"`
	Functions      []analyzer.FunctionAnalysis `json:"functions"`
	FileComplexity int                         `json:"file_complexity"`
}

// GenerateReport generates a JSON report from file analyses.
func GenerateReport(analyses []analyzer.FileAnalysis) []byte {
	report := buildReport(analyses)
	jsonReport, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return []byte("{}")
	}
	return jsonReport
}

func buildReport(analyses []analyzer.FileAnalysis) Report {
	var report Report
	report.TotalFiles = len(analyses)

	for _, fa := range analyses {
		fileReport := FileReport{FilePath: fa.FilePath}

		for _, fn := range fa.Functions {
			fileReport.Functions = append(fileReport.Functions, fn)
			report.TotalFunctions++
			report.TotalComplexity += fn.Cyclomatic
			fileReport.FileComplexity += fn.Cyclomatic
		}

		report.Files = append(report.Files, fileReport)
	}

	return report
}
