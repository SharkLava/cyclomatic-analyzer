# Cyclomatic Analyzer

Cyclomatic Analyzer is a Go-based tool designed to analyze C source code and calculate the cyclomatic complexity of functions within the code. It provides insights into the complexity of your C codebase, helping developers identify areas that may need refactoring or closer attention during code reviews.

## Features

- Recursively analyze C source files (.c and .h) in a given directory
- Calculate cyclomatic complexity for each function
- Identify decision points in the code
- Generate a JSON report with detailed analysis results
- Concurrent processing for improved performance on multi-core systems

## Usage

To use Cyclomatic Analyzer, navigate to the directory containing your C source files and run:

```
cyclomatic-analyzer -path /path/to/c/files -output report.json
```

Options:
- `-path`: Path to the directory containing C source files (default: current directory)
- `-output`: Name of the output JSON report file (default: report.json)

## Output

The tool generates a JSON report with the following structure:

```json
{
  "total_files": 10,
  "total_functions": 50,
  "total_complexity": 150,
  "files": [
    {
      "file_path": "/path/to/file.c",
      "functions": [
        {
          "name": "function_name",
          "cyclomatic": 5,
          "decision_points": ["if", "for", "while"]
        }
      ],
      "file_complexity": 15
    }
  ]
}
```

## Development

### Project Structure

- `cmd/analyzer/`: Contains the main application entry point
- `internal/`: Internal packages
  - `analyzer/`: Core analysis logic
  - `parser/`: C code parsing functionality
  - `report/`: Report generation
- `pkg/`: Shared utility functions

### Building

To build the project, run:

```
go build -o cyclomatic-analyzer ./cmd/analyzer
```
