package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

	// Define regex patterns to match struct definitions
	structPattern := regexp.MustCompile(`type\s+(\w+)\s+struct\s*{`)

	// Maps to store struct names found in *.pb.go files
	structsInPb := make(map[string]bool)
	// List to keep track of Go files to process
	files := make([]string, 0)

	// First pass: Gather struct names from *.pb.go files
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Add file to the list of files to process
		files = append(files, path)

		// Process only *.pb.go files to gather struct names
		if strings.HasSuffix(path, ".pb.go") {
			return processPbFile(path, structPattern, structsInPb)
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error walking through directory:", err)
		return
	}

	// Second pass: Remove duplicate structs from *.pb.go files
	for _, filePath := range files {
		if strings.HasSuffix(filePath, ".pb.go") {
			err = removeDuplicateStructs(filePath, structPattern, structsInPb)
			if err != nil {
				fmt.Printf("Error processing file %s: %v\n", filePath, err)
			}
		}
	}

	fmt.Println("Processing completed. Duplicate structs have been removed.")
}

// processPbFile extracts struct names from *.pb.go files and stores them in a map
func processPbFile(filePath string, structPattern *regexp.Regexp, structsInPb map[string]bool) error {
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
			structsInPb[structName] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// removeDuplicateStructs removes structs and their preceding comments that are duplicated in the package from *.pb.go files
func removeDuplicateStructs(filePath string, structPattern *regexp.Regexp, structsInPb map[string]bool) error {
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
	structName := ""
	commentBuffer := []string{}

	for scanner.Scan() {
		line := scanner.Text()

		// Check for struct definitions
		if structPattern.MatchString(line) {
			matches := structPattern.FindStringSubmatch(line)
			structName = matches[1]

			if structsInPb[structName] {
				inStruct = true
				commentBuffer = nil // Clear any comments preceding this struct
				// Skip the struct block to remove it from the file
				continue
			}
		}

		if inStruct {
			if strings.HasSuffix(line, "}") {
				inStruct = false
				// Skip the closing brace of the struct
				continue
			}
			// Skip lines inside the struct
			continue
		}

		// Handle single-line comments preceding structs
		if strings.HasPrefix(line, "//") {
			commentBuffer = append(commentBuffer, line)
		} else {
			// Write comments if not followed by a struct definition
			if len(commentBuffer) > 0 {
				for _, comment := range commentBuffer {
					_, _ = writer.WriteString(comment + "\n")
				}
				commentBuffer = nil
			}
			_, _ = writer.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Write remaining comments if any
	if len(commentBuffer) > 0 {
		for _, comment := range commentBuffer {
			_, _ = writer.WriteString(comment + "\n")
		}
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
