package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"io"
)

func main() {
	// Get the current working directory
	projectRoot, err := os.Getwd()
	if err != nil {
		fmt.Println("Error determining the project root directory:", err)
		return
	}

	// Define the directory to scan for Go files
	dir := filepath.Join(projectRoot, "pkg", "apis", "application", "v1alpha1")

	// Define regex patterns to match struct and function definitions
	structPattern := regexp.MustCompile(`type\s+(\w+)\s+struct\s*{`)
	funcPattern := regexp.MustCompile(`func\s+\((\w+)\s+\*?(\w+)\)\s+(\w+)\(`)

	// Maps to store struct and function names found in other Go files
	structsInOtherFiles := make(map[string]bool)
	functionsInOtherFiles := make(map[string]bool)

	// First pass: Gather struct and function names from non-pb.go files
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || strings.HasSuffix(path, ".pb.go") || !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Process non-pb.go files
		return collectStructAndFunctionNames(path, structPattern, funcPattern, structsInOtherFiles, functionsInOtherFiles)
	})

	if err != nil {
		fmt.Println("Error collecting structs and functions from non-pb.go files:", err)
		return
	}

	// Second pass: Process *.pb.go files and remove duplicates
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".pb.go") {
			return nil
		}

		// Process pb.go files and remove duplicates
		return removeDuplicateStructsAndFunctions(path, structPattern, funcPattern, structsInOtherFiles, functionsInOtherFiles)
	})

	if err != nil {
		fmt.Println("Error processing pb.go files:", err)
		return
	}

	fmt.Println("Processing completed. Duplicate structs and functions have been removed.")
}

// collectStructAndFunctionNames gathers struct and function names from non-pb.go files
func collectStructAndFunctionNames(filePath string, structPattern, funcPattern *regexp.Regexp, structsInOtherFiles, functionsInOtherFiles map[string]bool) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Check for struct definitions
		if structPattern.MatchString(line) {
			matches := structPattern.FindStringSubmatch(line)
			structName := matches[1]
			structsInOtherFiles[structName] = true
		}

		// Check for function definitions
		if funcPattern.MatchString(line) {
			matches := funcPattern.FindStringSubmatch(line)
			receiver := matches[1]
			receiverType := matches[2]
			funcName := matches[3]
			functionsInOtherFiles[receiver+"."+funcName] = true
			functionsInOtherFiles[receiverType+"."+funcName] = true // Also consider receiver type
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// removeDuplicateStructsAndFunctions removes duplicates from pb.go files
func removeDuplicateStructsAndFunctions(filePath string, structPattern, funcPattern *regexp.Regexp, structsInOtherFiles, functionsInOtherFiles map[string]bool) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a temporary file to write the modified content
	tempFile, err := os.CreateTemp(filepath.Dir(filePath), "temp_*.go")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(tempFile)

	inStruct := false
	inFunc := false
	bracesCount := 0 // Track the number of braces to handle nested curly braces

	for scanner.Scan() {
		line := scanner.Text()

		// Check for struct definitions
		if structPattern.MatchString(line) && !inStruct && !inFunc {
			matches := structPattern.FindStringSubmatch(line)
			structName := matches[1]

			if structsInOtherFiles[structName] {
				// Skip the entire struct block
				inStruct = true
				bracesCount = 1
				continue
			}
		}

		// Check for function definitions
		if funcPattern.MatchString(line) && !inStruct && !inFunc {
			matches := funcPattern.FindStringSubmatch(line)
			receiver := matches[1]
			receiverType := matches[2]
			funcName := matches[3]

			if functionsInOtherFiles[receiver+"."+funcName] || functionsInOtherFiles[receiverType+"."+funcName] {
				// Skip the entire function block
				inFunc = true
				bracesCount = 1
				continue
			}
		}

		// If inside a struct or function, track curly braces to know when to stop skipping
		if inStruct || inFunc {
			bracesCount += strings.Count(line, "{") - strings.Count(line, "}")
			if bracesCount == 0 {
				// End of the struct or function block
				inStruct = false
				inFunc = false
			}
			continue
		}

		// Write non-duplicate lines to the temporary file
		_, _ = writer.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	writer.Flush()

	// Replace the original file with the modified temporary file
	if err := replaceFile(filePath, tempFile.Name()); err != nil {
		return err
	}

	return nil
}

// replaceFile replaces the original file with the new file
func replaceFile(originalPath, newPath string) error {
	err := os.Rename(newPath, originalPath)
	if err != nil {
		// If renaming fails, copy and then remove the old file
		if err := copyFile(newPath, originalPath); err != nil {
			return err
		}
		if err := os.Remove(newPath); err != nil {
			return err
		}
	}
	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}
