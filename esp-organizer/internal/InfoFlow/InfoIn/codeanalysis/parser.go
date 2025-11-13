package codeanalysis

import (
	"encoding/json"
	"esp-organizer/internal/InfoFlow/InfoStore/db"
	"io"
	"net/http"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type CodeFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  []Parameter            `json:"parameters"`
	ReturnType  string                 `json:"return_type,omitempty"`
	Purpose     string                 `json:"purpose"`
	Language    string                 `json:"language"`
	FilePath    string                 `json:"file_path"`
	LineStart   int                    `json:"line_start"`
	LineEnd     int                    `json:"line_end"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type Parameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type CodeAnalyzer struct {
	supportedExtensions map[string]string
}

func NewCodeAnalyzer() *CodeAnalyzer {
	return &CodeAnalyzer{
		supportedExtensions: map[string]string{
			".go":  "go",
			".js":  "javascript",
			".jsx": "react",
			".ts":  "typescript",
			".tsx": "typescript-react",
			".php": "php",
			".py":  "python",
		},
	}
}

// AnalyzeReactFile parses React/JavaScript files
func (ca *CodeAnalyzer) AnalyzeReactFile(filePath string, content string) ([]CodeFunction, error) {
	functions := []CodeFunction{}

	// React component detection patterns
	componentPattern := regexp.MustCompile(`(?m)^(?:export\s+)?(?:const|function|class)\s+(\w+)\s*(?:\(.*?\)|=\s*\(.*?\)\s*=>|extends\s+\w+)`)

	// Find React components
	matches := componentPattern.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			componentName := match[1]

			// Extract component logic
			purpose := extractReactComponentPurpose(content, componentName)
			props := extractReactProps(content, componentName)
			hooks := extractReactHooks(content)

			function := CodeFunction{
				Name:        componentName,
				Description: generateReactDescription(componentName, purpose, props, hooks),
				Purpose:     purpose,
				Language:    "react",
				FilePath:    filePath,
				Parameters:  props,
				Metadata: map[string]interface{}{
					"component_type": detectComponentType(content, componentName),
					"hooks_used":     hooks,
					"imports":        extractImports(content),
					"exports":        extractExports(content),
				},
			}

			functions = append(functions, function)
		}
	}

	return functions, nil
}

// AnalyzePHPFile parses PHP files
func (ca *CodeAnalyzer) AnalyzePHPFile(filePath string, content string) ([]CodeFunction, error) {
	functions := []CodeFunction{}

	// PHP function detection patterns
	functionPattern := regexp.MustCompile(`(?m)^\s*(?:public|private|protected|static)?\s*function\s+(\w+)\s*\(([^)]*)\)`)

	// Find PHP functions
	matches := functionPattern.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 2 {
			functionName := match[1]
			paramString := match[2]

			parameters := parsePHPParameters(paramString)
			purpose := extractPHPFunctionPurpose(content, functionName)

			function := CodeFunction{
				Name:        functionName,
				Description: generatePHPDescription(functionName, purpose, parameters),
				Purpose:     purpose,
				Language:    "php",
				FilePath:    filePath,
				Parameters:  parameters,
				Metadata: map[string]interface{}{
					"visibility": extractPHPVisibility(content, functionName),
					"class_name": extractPHPClassName(content),
					"namespace":  extractPHPNamespace(content),
				},
			}

			functions = append(functions, function)
		}
	}

	return functions, nil
}

// Helper functions
func extractReactComponentPurpose(content string, componentName string) string {
	// Look for JSDoc comments or component descriptions
	commentPattern := regexp.MustCompile(`/\*\*([^*]|\*[^/])*\*/\s*(?:export\s+)?(?:const|function|class)\s+` + componentName)
	if match := commentPattern.FindString(content); match != "" {
		return cleanComment(match)
	}

	// Fallback: analyze JSX structure
	return analyzeJSXStructure(content, componentName)
}

func extractReactProps(content string, componentName string) []Parameter {
	propsPattern := regexp.MustCompile(`(?:function\s+` + componentName + `\s*\(\s*\{([^}]+)\}|const\s+` + componentName + `\s*=\s*\(\s*\{([^}]+)\})`)
	match := propsPattern.FindStringSubmatch(content)

	if len(match) > 1 {
		propsStr := match[1]
		if propsStr == "" && len(match) > 2 {
			propsStr = match[2]
		}
		return parseJSParameters(propsStr)
	}

	return []Parameter{}
}

func generateReactDescription(name string, purpose string, props []Parameter, hooks []string) string {
	desc := "React component: " + name + ". "

	if purpose != "" {
		desc += purpose + ". "
	}

	if len(props) > 0 {
		propNames := make([]string, len(props))
		for i, prop := range props {
			propNames[i] = prop.Name
		}
		desc += "Accepts props: " + strings.Join(propNames, ", ") + ". "
	}

	if len(hooks) > 0 {
		desc += "Uses React hooks: " + strings.Join(hooks, ", ") + ". "
	}

	return desc
}

// Add missing helper functions that are referenced but not implemented
func cleanComment(comment string) string {
	// Remove comment markers and clean up
	comment = strings.ReplaceAll(comment, "/**", "")
	comment = strings.ReplaceAll(comment, "*/", "")
	comment = strings.ReplaceAll(comment, "*", "")
	return strings.TrimSpace(comment)
}

func analyzeJSXStructure(content string, componentName string) string {
	// Simple analysis based on JSX elements
	if strings.Contains(content, "<div") {
		return "Renders a div-based layout"
	}
	if strings.Contains(content, "<form") {
		return "Renders a form component"
	}
	if strings.Contains(content, "<button") {
		return "Contains interactive button elements"
	}
	return "React functional component"
}

func extractReactHooks(content string) []string {
	hooks := []string{}
	hookPatterns := map[string]string{
		`useState`:    "useState",
		`useEffect`:   "useEffect",
		`useContext`:  "useContext",
		`useReducer`:  "useReducer",
		`useCallback`: "useCallback",
		`useMemo`:     "useMemo",
	}

	for pattern, hook := range hookPatterns {
		if strings.Contains(content, pattern) {
			hooks = append(hooks, hook)
		}
	}

	return hooks
}

func detectComponentType(content string, componentName string) string {
	if strings.Contains(content, "class "+componentName) {
		return "class_component"
	}
	if strings.Contains(content, "const "+componentName) {
		return "functional_component"
	}
	return "function_component"
}

func extractImports(content string) []string {
	importPattern := regexp.MustCompile(`import\s+.*?\s+from\s+['"]([^'"]+)['"]`)
	matches := importPattern.FindAllStringSubmatch(content, -1)

	imports := []string{}
	for _, match := range matches {
		if len(match) > 1 {
			imports = append(imports, match[1])
		}
	}

	return imports
}

func extractExports(content string) []string {
	exportPattern := regexp.MustCompile(`export\s+(?:default\s+)?(\w+)`)
	matches := exportPattern.FindAllStringSubmatch(content, -1)

	exports := []string{}
	for _, match := range matches {
		if len(match) > 1 {
			exports = append(exports, match[1])
		}
	}

	return exports
}

func parseJSParameters(propsStr string) []Parameter {
	params := []Parameter{}

	// Simple parsing - split by comma and clean up
	parts := strings.Split(propsStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			// Handle destructured props like "title, onClick, disabled"
			params = append(params, Parameter{
				Name: part,
				Type: "any", // JavaScript doesn't have explicit types
			})
		}
	}

	return params
}

func parsePHPParameters(paramString string) []Parameter {
	params := []Parameter{}

	if paramString == "" {
		return params
	}

	// Split by comma and parse each parameter
	parts := strings.Split(paramString, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			// Handle PHP parameters like "$name", "$value = null", "string $title"
			name := part
			paramType := "mixed"

			// Extract type hints
			if strings.Contains(part, " ") {
				typeParts := strings.Fields(part)
				if len(typeParts) >= 2 {
					paramType = typeParts[0]
					name = typeParts[1]
				}
			}

			// Remove $ prefix and default values
			name = strings.Split(name, "=")[0]
			name = strings.TrimPrefix(strings.TrimSpace(name), "$")

			params = append(params, Parameter{
				Name: name,
				Type: paramType,
			})
		}
	}

	return params
}

func extractPHPFunctionPurpose(content string, functionName string) string {
	// Look for PHPDoc comments before the function
	commentPattern := regexp.MustCompile(`/\*\*([^*]|\*[^/])*\*/\s*(?:public|private|protected|static)?\s*function\s+` + functionName)
	if match := commentPattern.FindString(content); match != "" {
		return cleanComment(match)
	}

	return "PHP function: " + functionName
}

func generatePHPDescription(functionName string, purpose string, parameters []Parameter) string {
	desc := "PHP function: " + functionName + ". "

	if purpose != "" {
		desc += purpose + ". "
	}

	if len(parameters) > 0 {
		paramNames := make([]string, len(parameters))
		for i, param := range parameters {
			paramNames[i] = param.Name + " (" + param.Type + ")"
		}
		desc += "Parameters: " + strings.Join(paramNames, ", ") + ". "
	}

	return desc
}

func extractPHPVisibility(content string, functionName string) string {
	visibilityPattern := regexp.MustCompile(`(public|private|protected)\s+function\s+` + functionName)
	match := visibilityPattern.FindStringSubmatch(content)

	if len(match) > 1 {
		return match[1]
	}

	return "public" // Default visibility
}

func extractPHPClassName(content string) string {
	classPattern := regexp.MustCompile(`class\s+(\w+)`)
	match := classPattern.FindStringSubmatch(content)

	if len(match) > 1 {
		return match[1]
	}

	return ""
}

func extractPHPNamespace(content string) string {
	namespacePattern := regexp.MustCompile(`namespace\s+([^;]+);`)
	match := namespacePattern.FindStringSubmatch(content)

	if len(match) > 1 {
		return match[1]
	}

	return ""
}

// CodeAnalysisUploadHandler handles code file uploads for analysis and returns completion status
func CodeAnalysisUploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Parse uploaded file
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file in chunks to simulate progress (for demonstration)
	const bufSize = 1024 * 1024 // 1MB
	var totalRead int64
	buf := make([]byte, bufSize)
	for {
		n, err := file.Read(buf)
		if n > 0 {
			totalRead += int64(n)
			// Here you could update progress in a store or log
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
	}

	// Respond when upload is complete
	response := map[string]interface{}{
		"status":   "uploaded",
		"filename": header.Filename,
		"size":     totalRead,
		"message":  "File uploaded and ready for analysis.",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetCodeFunctionsHandler returns processed code functions
func GetCodeFunctionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	mongOb, err := db.NewFromEnv()
	if err != nil {
		http.Error(w, "Failed to connect to MongoDB: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer mongOb.Client.Disconnect(r.Context())
	codeFunctionsCollection := mongOb.Database.Collection("code_functions")

	cursor, err := codeFunctionsCollection.Find(r.Context(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to query code functions", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	var functions []map[string]interface{}
	for cursor.Next(r.Context()) {
		var function bson.M
		if err := cursor.Decode(&function); err != nil {
			continue
		}
		functions = append(functions, function)
	}

	response := map[string]interface{}{
		"code_functions": functions,
		"total_count":    len(functions),
	}

	json.NewEncoder(w).Encode(response)
}
