package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/SharkLava/cyclomatic-analyzer/internal/analyzer"
	"github.com/SharkLava/cyclomatic-analyzer/internal/parser"
	"github.com/SharkLava/cyclomatic-analyzer/internal/report"
	"github.com/SharkLava/cyclomatic-analyzer/pkg/utils"
)

func main() {
	path, output := parseFlags()
	cFiles := findCFiles(path)

	if len(cFiles) == 0 {
		log.Println("No C source files found.")
		os.Exit(0)
	}

	analyses := analyzeFiles(cFiles)
	writeReport(analyses, output)
}

func parseFlags() (string, string) {
	path := flag.String("path", ".", "Path to the directory containing C source files")
	output := flag.String("output", "report.json", "Output report file (JSON format)")
	flag.Parse()

	absPath, err := filepath.Abs(*path)
	if err != nil {
		log.Fatalf("Invalid path: %v", err)
	}

	return absPath, *output
}

func findCFiles(path string) []string {
	cFiles, err := utils.FindCFiles(path)
	if err != nil {
		log.Fatalf("Error finding C files: %v", err)
	}
	return cFiles
}

func analyzeFiles(cFiles []string) []analyzer.FileAnalysis {
	numWorkers := runtime.NumCPU()
	fileChan := make(chan string, numWorkers)
	resultChan := make(chan analyzer.FileAnalysis, numWorkers)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(fileChan, resultChan, &wg)
	}

	go func() {
		for _, file := range cFiles {
			fileChan <- file
		}
		close(fileChan)
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var analyses []analyzer.FileAnalysis
	for analysis := range resultChan {
		analyses = append(analyses, analysis)
	}

	return analyses
}

func worker(fileChan <-chan string, resultChan chan<- analyzer.FileAnalysis, wg *sync.WaitGroup) {
	defer wg.Done()
	for filePath := range fileChan {
		code, err := utils.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading file %s: %v", filePath, err)
			continue
		}

		parsedFunctions, err := parser.ParseFunctions(code)
		if err != nil {
			log.Printf("Error parsing file %s: %v", filePath, err)
			continue
		}

		// Convert parser.Function to analyzer.Function
		functions := make([]analyzer.Function, len(parsedFunctions))
		for i, pf := range parsedFunctions {
			functions[i] = analyzer.Function{Name: pf.Name, Body: pf.Body}
		}

		funcAnalyses := analyzer.AnalyzeFunctions(functions)
		fileAnalysis := analyzer.FileAnalysis{
			FilePath:  filePath,
			Functions: funcAnalyses,
		}

		resultChan <- fileAnalysis
	}
}

func writeReport(analyses []analyzer.FileAnalysis, output string) {
	reportData := report.GenerateReport(analyses)
	err := utils.WriteFile(output, reportData)
	if err != nil {
		log.Fatalf("Error writing report: %v", err)
	}

	fmt.Printf("Analysis complete. Report saved to %s\n", output)
}
