package main

import (
	"os"
	"path/filepath"
	"strings"
)

// CodeFile represents a source code file with its path and contents.
type CodeFile struct {
	Path    string
	Content string
}

// targetExtensions is a slice containing the file extensions we want to scan.
// In Go, a slice is a dynamically sized, flexible view into an array.
var targetExtensions = []string{
	".go",   // Go
	".js",   // JavaScript
	".ts",   // TypeScript
	".py",   // Python
	".java", // Java
	".rs",   // Rust
	".cpp",  // C++
	".c",    // C
	".h",    // Header files
}

// scanDirectory walks the directory at the given root path, reads all matched
// source files, and returns them as a slice of CodeFile.
func scanDirectory(root string) ([]CodeFile, error) {
	var files []CodeFile

	// filepath.WalkDir recursively walks the file tree rooted at 'root'.
	// It calls the provided function for each file and directory it encounters.
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		// 1. If there's an error accessing the path, return it so WalkDir can handle or stop
		if err != nil {
			return err
		}

		// 2. Check if the current item is a directory
		if d.IsDir() {
			name := d.Name()
			// Skip directories we want to ignore (like version control or dependencies)
			if name == ".git" || name == "node_modules" || name == ".vscode" || name == "dist" || name == "bin" {
				// filepath.SkipDir is a special sentinel error. Returning it tells WalkDir:
				// "Do not enter this directory, skip everything inside it and continue."
				return filepath.SkipDir
			}
			return nil // Continue walking other directories
		}

		// 3. Check if the file matches our target extensions
		ext := filepath.Ext(path)
		if isTargetExtension(ext) {
			// os.ReadFile reads the entire file into memory as a slice of bytes ([]byte).
			bytes, err := os.ReadFile(path)
			if err != nil {
				// If reading a single file fails, we print a message but continue scanning others
				return nil
			}

			// We append the successfully read file to our list
			// string(bytes) converts the slice of bytes into a Go string.
			files = append(files, CodeFile{
				Path:    path,
				Content: string(bytes),
			})
		}

		return nil // Return nil to signal WalkDir to continue to the next item
	})

	// If the entire walking process failed (e.g. root dir doesn't exist), return the error.
	if err != nil {
		return nil, err
	}

	return files, nil
}

// isTargetExtension checks if a given extension matches our target list.
func isTargetExtension(ext string) bool {
	// Convert extension to lowercase to handle cases like .PY or .JS
	ext = strings.ToLower(ext)
	for _, target := range targetExtensions {
		if ext == target {
			return true
		}
	}
	return false
}