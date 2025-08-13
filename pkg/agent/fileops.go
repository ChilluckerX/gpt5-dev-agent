package agent

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// FileOperations handles file access and operations for the agent
type FileOperations struct {
	workingDir string
	allowedExts []string
	maxFileSize int64
}

// NewFileOperations creates a new file operations handler
func NewFileOperations() *FileOperations {
	workingDir, _ := os.Getwd()
	return &FileOperations{
		workingDir: workingDir,
		allowedExts: []string{
			".go", ".py", ".js", ".ts", ".java", ".rs", ".cpp", ".c", ".h",
			".md", ".txt", ".json", ".yaml", ".yml", ".toml", ".xml",
			".html", ".css", ".sql", ".sh", ".bat", ".dockerfile",
			".gitignore", ".env", "makefile",
		},
		maxFileSize: 10 * 1024 * 1024, // 10MB limit
	}
}

// ReadFile reads a specific file and returns its content
func (fo *FileOperations) ReadFile(filename string) (string, error) {
	// Security check: ensure file is within working directory
	fullPath := filepath.Join(fo.workingDir, filename)
	if !strings.HasPrefix(fullPath, fo.workingDir) {
		return "", fmt.Errorf("access denied: file outside working directory")
	}

	// Check if file exists
	info, err := os.Stat(fullPath)
	if err != nil {
		return "", fmt.Errorf("file not found: %s", filename)
	}

	// Check file size
	if info.Size() > fo.maxFileSize {
		return "", fmt.Errorf("file too large: %s (max %d bytes)", filename, fo.maxFileSize)
	}

	// Check if file extension is allowed
	ext := strings.ToLower(filepath.Ext(filename))
	if !fo.isAllowedExtension(ext) && !fo.isSpecialFile(filename) {
		return "", fmt.Errorf("file type not allowed: %s", ext)
	}

	// Read file content
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	return string(content), nil
}

// ListFiles lists all files in the current directory or specified path
func (fo *FileOperations) ListFiles(path string) ([]FileInfo, error) {
	var targetPath string
	if path == "" || path == "." {
		targetPath = fo.workingDir
	} else {
		targetPath = filepath.Join(fo.workingDir, path)
		// Security check
		if !strings.HasPrefix(targetPath, fo.workingDir) {
			return nil, fmt.Errorf("access denied: path outside working directory")
		}
	}

	var files []FileInfo
	err := filepath.WalkDir(targetPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files and directories (except important ones)
		name := d.Name()
		if strings.HasPrefix(name, ".") && !fo.isImportantHiddenFile(name) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip common ignore patterns
		if fo.shouldSkip(name) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return nil // Skip files we can't stat
			}

			// Get relative path from working directory
			relPath, err := filepath.Rel(fo.workingDir, path)
			if err != nil {
				relPath = path
			}

			files = append(files, FileInfo{
				Name:      name,
				Path:      relPath,
				Extension: strings.ToLower(filepath.Ext(name)),
				Category:  fo.categorizeFile(name),
				Size:      info.Size(),
				ModTime:   info.ModTime(),
			})
		}

		return nil
	})

	return files, err
}

// SearchFiles searches for files matching a pattern
func (fo *FileOperations) SearchFiles(pattern string) ([]FileInfo, error) {
	allFiles, err := fo.ListFiles("")
	if err != nil {
		return nil, err
	}

	var matches []FileInfo
	pattern = strings.ToLower(pattern)

	for _, file := range allFiles {
		// Search in filename
		if strings.Contains(strings.ToLower(file.Name), pattern) {
			matches = append(matches, file)
			continue
		}

		// Search in file path
		if strings.Contains(strings.ToLower(file.Path), pattern) {
			matches = append(matches, file)
		}
	}

	return matches, nil
}

// ReadMultipleFiles reads multiple files and returns their content
func (fo *FileOperations) ReadMultipleFiles(filenames []string) (map[string]string, error) {
	results := make(map[string]string)
	
	for _, filename := range filenames {
		content, err := fo.ReadFile(filename)
		if err != nil {
			results[filename] = fmt.Sprintf("Error reading file: %v", err)
		} else {
			results[filename] = content
		}
	}

	return results, nil
}

// GetFileTree returns a tree structure of the project
func (fo *FileOperations) GetFileTree(maxDepth int) (string, error) {
	var tree strings.Builder
	
	err := fo.buildTree(fo.workingDir, "", 0, maxDepth, &tree)
	if err != nil {
		return "", err
	}

	return tree.String(), nil
}

// buildTree recursively builds the file tree
func (fo *FileOperations) buildTree(dir, prefix string, depth, maxDepth int, tree *strings.Builder) error {
	if maxDepth > 0 && depth >= maxDepth {
		return nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for i, entry := range entries {
		name := entry.Name()
		
		// Skip hidden and ignored files
		if strings.HasPrefix(name, ".") && !fo.isImportantHiddenFile(name) {
			continue
		}
		if fo.shouldSkip(name) {
			continue
		}

		isLast := i == len(entries)-1
		var connector, newPrefix string
		
		if isLast {
			connector = "└── "
			newPrefix = prefix + "    "
		} else {
			connector = "├── "
			newPrefix = prefix + "│   "
		}

		tree.WriteString(prefix + connector + name)
		
		if entry.IsDir() {
			tree.WriteString("/\n")
			subDir := filepath.Join(dir, name)
			fo.buildTree(subDir, newPrefix, depth+1, maxDepth, tree)
		} else {
			// Add file size info
			if info, err := entry.Info(); err == nil {
				tree.WriteString(fmt.Sprintf(" (%s)", fo.formatFileSize(info.Size())))
			}
			tree.WriteString("\n")
		}
	}

	return nil
}

// Helper functions

func (fo *FileOperations) isAllowedExtension(ext string) bool {
	for _, allowed := range fo.allowedExts {
		if ext == allowed {
			return true
		}
	}
	return false
}

func (fo *FileOperations) isSpecialFile(filename string) bool {
	specialFiles := []string{
		"Makefile", "Dockerfile", "README", "LICENSE", "CHANGELOG",
		"go.mod", "go.sum", "package.json", "requirements.txt",
	}
	
	name := strings.ToLower(filepath.Base(filename))
	for _, special := range specialFiles {
		if strings.Contains(name, strings.ToLower(special)) {
			return true
		}
	}
	return false
}

func (fo *FileOperations) isImportantHiddenFile(name string) bool {
	important := []string{
		".env", ".gitignore", ".dockerignore", ".editorconfig",
		".eslintrc", ".prettierrc", ".babelrc", ".nvmrc",
	}
	
	for _, imp := range important {
		if name == imp {
			return true
		}
	}
	return false
}

func (fo *FileOperations) shouldSkip(name string) bool {
	skipPatterns := []string{
		"node_modules", "vendor", "target", "build", "dist",
		".git", ".svn", ".hg", "__pycache__", ".pytest_cache",
		"*.exe", "*.dll", "*.so", "*.dylib", "*.a", "*.o",
	}
	
	lowerName := strings.ToLower(name)
	for _, pattern := range skipPatterns {
		if strings.Contains(lowerName, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

func (fo *FileOperations) categorizeFile(filename string) FileCategory {
	name := strings.ToLower(filename)
	ext := strings.ToLower(filepath.Ext(filename))
	
	// Config files
	configFiles := []string{
		"go.mod", "go.sum", "package.json", "package-lock.json",
		"requirements.txt", "setup.py", "pyproject.toml",
		"cargo.toml", "cargo.lock", "pom.xml", "build.gradle",
		"dockerfile", "docker-compose.yml", ".gitignore",
		"makefile", "cmake", "config.json", "config.yaml",
	}
	for _, config := range configFiles {
		if name == config || strings.Contains(name, config) {
			return ConfigFile
		}
	}
	
	// Code files
	codeExts := []string{
		".go", ".py", ".js", ".ts", ".java", ".rs", ".cpp", ".c", ".h",
		".cs", ".php", ".rb", ".swift", ".kt", ".scala", ".clj",
	}
	for _, codeExt := range codeExts {
		if ext == codeExt {
			return CodeFile
		}
	}
	
	// Test files
	if strings.Contains(name, "test") || strings.Contains(name, "spec") {
		return TestFile
	}
	
	// Documentation
	docExts := []string{".md", ".txt", ".rst", ".adoc"}
	for _, docExt := range docExts {
		if ext == docExt {
			return DocumentFile
		}
	}
	
	// Build files
	buildFiles := []string{"makefile", "build.sh", "build.bat", "webpack.config.js"}
	for _, build := range buildFiles {
		if strings.Contains(name, build) {
			return BuildFile
		}
	}
	
	return UnknownFile
}

func (fo *FileOperations) formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}