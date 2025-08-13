package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ProjectContext handles project analysis and context management
type ProjectContext struct {
	currentDir    string
	projectName   string
	projectType   string
	files         []FileInfo
	directories   []string
	lastAnalyzed  time.Time
	analysis      ProjectAnalysis
}

// FileInfo represents information about a file
type FileInfo struct {
	Name      string
	Path      string
	Extension string
	Category  FileCategory
	Size      int64
	ModTime   time.Time
}

// FileCategory represents the category of a file
type FileCategory string

const (
	ConfigFile     FileCategory = "config"
	CodeFile       FileCategory = "code"
	DocumentFile   FileCategory = "document"
	DataFile       FileCategory = "data"
	BuildFile      FileCategory = "build"
	TestFile       FileCategory = "test"
	UnknownFile    FileCategory = "unknown"
)

// ProjectAnalysis contains the results of project analysis
type ProjectAnalysis struct {
	ProjectType    string
	Technologies   []string
	Dependencies   []string
	Structure      ProjectStructure
	Insights       []string
}

// ProjectStructure represents the project's structure
type ProjectStructure struct {
	HasTests       bool
	HasDocs        bool
	HasConfig      bool
	HasBuild       bool
	PackageManager string
	MainFiles      []string
}

// NewProjectContext creates a new project context
func NewProjectContext() *ProjectContext {
	currentDir, _ := os.Getwd()
	projectName := filepath.Base(currentDir)
	
	ctx := &ProjectContext{
		currentDir:  currentDir,
		projectName: projectName,
	}
	
	ctx.Refresh()
	return ctx
}

// Refresh re-analyzes the project
func (pc *ProjectContext) Refresh() error {
	pc.lastAnalyzed = time.Now()
	
	// Analyze files and directories
	if err := pc.analyzeStructure(); err != nil {
		return err
	}
	
	// Detect project type and technologies
	pc.detectProjectType()
	pc.detectTechnologies()
	pc.generateInsights()
	
	return nil
}

// analyzeStructure analyzes the project's file structure
func (pc *ProjectContext) analyzeStructure() error {
	entries, err := os.ReadDir(pc.currentDir)
	if err != nil {
		return err
	}
	
	pc.files = []FileInfo{}
	pc.directories = []string{}
	
	for _, entry := range entries {
		name := entry.Name()
		
		// Skip hidden files (except important ones)
		if strings.HasPrefix(name, ".") && !pc.isImportantHiddenFile(name) {
			continue
		}
		
		if entry.IsDir() {
			pc.directories = append(pc.directories, name)
		} else {
			info, err := entry.Info()
			if err != nil {
				continue
			}
			
			fileInfo := FileInfo{
				Name:      name,
				Extension: strings.ToLower(filepath.Ext(name)),
				Category:  pc.categorizeFile(name),
				Size:      info.Size(),
				ModTime:   info.ModTime(),
			}
			pc.files = append(pc.files, fileInfo)
		}
	}
	
	return nil
}

// isImportantHiddenFile checks if a hidden file is important for analysis
func (pc *ProjectContext) isImportantHiddenFile(name string) bool {
	importantFiles := []string{
		".env", ".gitignore", ".dockerignore", ".editorconfig",
		".eslintrc", ".prettierrc", ".babelrc", ".nvmrc",
	}
	
	for _, important := range importantFiles {
		if name == important {
			return true
		}
	}
	return false
}

// categorizeFile determines the category of a file
func (pc *ProjectContext) categorizeFile(filename string) FileCategory {
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

// detectProjectType determines the primary project type
func (pc *ProjectContext) detectProjectType() {
	// Check for specific project indicators
	for _, file := range pc.files {
		switch file.Name {
		case "go.mod":
			pc.projectType = "Go"
			return
		case "package.json":
			pc.projectType = "Node.js/JavaScript"
			return
		case "requirements.txt", "setup.py", "pyproject.toml":
			pc.projectType = "Python"
			return
		case "Cargo.toml":
			pc.projectType = "Rust"
			return
		case "pom.xml":
			pc.projectType = "Java/Maven"
			return
		case "build.gradle":
			pc.projectType = "Java/Gradle"
			return
		case "Dockerfile":
			if pc.projectType == "" {
				pc.projectType = "Docker"
			}
		}
	}
	
	// Check by file extensions if no specific indicators found
	if pc.projectType == "" {
		extCounts := make(map[string]int)
		for _, file := range pc.files {
			if file.Category == CodeFile {
				extCounts[file.Extension]++
			}
		}
		
		maxCount := 0
		var primaryExt string
		for ext, count := range extCounts {
			if count > maxCount {
				maxCount = count
				primaryExt = ext
			}
		}
		
		switch primaryExt {
		case ".go":
			pc.projectType = "Go"
		case ".py":
			pc.projectType = "Python"
		case ".js", ".ts":
			pc.projectType = "JavaScript/TypeScript"
		case ".java":
			pc.projectType = "Java"
		case ".rs":
			pc.projectType = "Rust"
		case ".cpp", ".c":
			pc.projectType = "C/C++"
		default:
			pc.projectType = "Mixed/Unknown"
		}
	}
}

// detectTechnologies identifies technologies used in the project
func (pc *ProjectContext) detectTechnologies() {
	pc.analysis.Technologies = []string{}
	
	// Add primary language
	if pc.projectType != "" && pc.projectType != "Mixed/Unknown" {
		pc.analysis.Technologies = append(pc.analysis.Technologies, pc.projectType)
	}
	
	// Check for specific technologies
	for _, file := range pc.files {
		switch file.Name {
		case "docker-compose.yml":
			pc.addTechnology("Docker Compose")
		case "Dockerfile":
			pc.addTechnology("Docker")
		case "webpack.config.js":
			pc.addTechnology("Webpack")
		case ".eslintrc", ".eslintrc.json":
			pc.addTechnology("ESLint")
		case "jest.config.js":
			pc.addTechnology("Jest")
		case "pytest.ini":
			pc.addTechnology("Pytest")
		}
	}
	
	// Check directories for frameworks
	for _, dir := range pc.directories {
		switch dir {
		case "node_modules":
			pc.addTechnology("Node.js")
		case "vendor":
			pc.addTechnology("Go Modules")
		case "venv", "env", ".venv":
			pc.addTechnology("Python Virtual Environment")
		case "target":
			if pc.projectType == "Rust" {
				pc.addTechnology("Cargo")
			} else if pc.projectType == "Java/Maven" {
				pc.addTechnology("Maven")
			}
		}
	}
}

// addTechnology adds a technology if not already present
func (pc *ProjectContext) addTechnology(tech string) {
	for _, existing := range pc.analysis.Technologies {
		if existing == tech {
			return
		}
	}
	pc.analysis.Technologies = append(pc.analysis.Technologies, tech)
}

// generateInsights creates insights about the project
func (pc *ProjectContext) generateInsights() {
	pc.analysis.Insights = []string{}
	
	// Analyze project structure
	pc.analysis.Structure = ProjectStructure{}
	
	for _, file := range pc.files {
		switch file.Category {
		case TestFile:
			pc.analysis.Structure.HasTests = true
		case DocumentFile:
			pc.analysis.Structure.HasDocs = true
		case ConfigFile:
			pc.analysis.Structure.HasConfig = true
		case BuildFile:
			pc.analysis.Structure.HasBuild = true
		}
		
		if file.Name == "main.go" || file.Name == "main.py" || file.Name == "index.js" {
			pc.analysis.Structure.MainFiles = append(pc.analysis.Structure.MainFiles, file.Name)
		}
	}
	
	// Generate insights based on analysis
	if pc.analysis.Structure.HasTests {
		pc.analysis.Insights = append(pc.analysis.Insights, "Project includes test files")
	}
	if pc.analysis.Structure.HasDocs {
		pc.analysis.Insights = append(pc.analysis.Insights, "Project has documentation")
	}
	if len(pc.analysis.Structure.MainFiles) > 0 {
		pc.analysis.Insights = append(pc.analysis.Insights, fmt.Sprintf("Entry points: %s", strings.Join(pc.analysis.Structure.MainFiles, ", ")))
	}
	
	// Check for common patterns
	if len(pc.directories) > 5 {
		pc.analysis.Insights = append(pc.analysis.Insights, "Well-organized project structure")
	}
	
	codeFileCount := 0
	for _, file := range pc.files {
		if file.Category == CodeFile {
			codeFileCount++
		}
	}
	
	if codeFileCount > 10 {
		pc.analysis.Insights = append(pc.analysis.Insights, "Large codebase with multiple modules")
	} else if codeFileCount > 3 {
		pc.analysis.Insights = append(pc.analysis.Insights, "Medium-sized project")
	} else {
		pc.analysis.Insights = append(pc.analysis.Insights, "Small/focused project")
	}
}

// GetProjectInfo returns a formatted string with project information
func (pc *ProjectContext) GetProjectInfo() string {
	var info strings.Builder
	
	info.WriteString(fmt.Sprintf("Project: %s\n", pc.projectName))
	info.WriteString(fmt.Sprintf("Type: %s\n", pc.projectType))
	
	if len(pc.analysis.Technologies) > 0 {
		info.WriteString(fmt.Sprintf("Technologies: %s\n", strings.Join(pc.analysis.Technologies, ", ")))
	}
	
	// File summary
	configFiles := []string{}
	codeFiles := []string{}
	for _, file := range pc.files {
		if file.Category == ConfigFile {
			configFiles = append(configFiles, file.Name)
		} else if file.Category == CodeFile {
			codeFiles = append(codeFiles, file.Name)
		}
	}
	
	if len(configFiles) > 0 {
		info.WriteString(fmt.Sprintf("Config files: %s\n", strings.Join(configFiles, ", ")))
	}
	if len(codeFiles) > 0 {
		info.WriteString(fmt.Sprintf("Code files: %s\n", strings.Join(codeFiles, ", ")))
	}
	if len(pc.directories) > 0 {
		info.WriteString(fmt.Sprintf("Directories: %s\n", strings.Join(pc.directories, ", ")))
	}
	
	if len(pc.analysis.Insights) > 0 {
		info.WriteString(fmt.Sprintf("Insights: %s\n", strings.Join(pc.analysis.Insights, ", ")))
	}
	
	return info.String()
}

// GetCurrentDir returns the current working directory
func (pc *ProjectContext) GetCurrentDir() string {
	return pc.currentDir
}

// EnhanceMessage adds project context to a user message
func (pc *ProjectContext) EnhanceMessage(message string) string {
	// For now, just return the original message
	// This could be enhanced to add relevant context based on the message content
	return message
}

// GetProjectType returns the detected project type
func (pc *ProjectContext) GetProjectType() string {
	return pc.projectType
}

// GetAnalysis returns the complete project analysis
func (pc *ProjectContext) GetAnalysis() ProjectAnalysis {
	return pc.analysis
}