package templates

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

// Engine handles template rendering
type Engine struct {
	templates map[string]*template.Template
	baseDir   string
	funcMap   template.FuncMap
}

// New creates a new template engine
func New() *Engine {
	return &Engine{
		templates: make(map[string]*template.Template),
		baseDir:   "templates",
		funcMap:   defaultFuncMap(),
	}
}

// SetBaseDir sets the base directory for templates
func (e *Engine) SetBaseDir(dir string) {
	e.baseDir = dir
}

// AddFunc adds a function to the template function map
func (e *Engine) AddFunc(name string, fn interface{}) {
	if e.funcMap == nil {
		e.funcMap = make(template.FuncMap)
	}
	e.funcMap[name] = fn
}

// LoadTemplates loads all templates from the base directory
func (e *Engine) LoadTemplates() error {
	if e.baseDir == "" {
		return nil
	}
	
	pattern := filepath.Join(e.baseDir, "*.html")
	templates, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to load templates: %v", err)
	}
	
	for _, templateFile := range templates {
		name := strings.TrimSuffix(filepath.Base(templateFile), ".html")
		
		tmpl, err := template.New(name).Funcs(e.funcMap).ParseFiles(templateFile)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %v", templateFile, err)
		}
		
		e.templates[name] = tmpl
	}
	
	return nil
}

// Render renders a template with data
func (e *Engine) Render(w io.Writer, name string, data interface{}) error {
	tmpl, exists := e.templates[name]
	if !exists {
		// Try to load the template dynamically
		if err := e.loadTemplate(name); err != nil {
			return fmt.Errorf("template %s not found: %v", name, err)
		}
		tmpl = e.templates[name]
	}
	
	return tmpl.Execute(w, data)
}

// loadTemplate loads a single template
func (e *Engine) loadTemplate(name string) error {
	templateFile := filepath.Join(e.baseDir, name+".html")
	
	tmpl, err := template.New(name).Funcs(e.funcMap).ParseFiles(templateFile)
	if err != nil {
		return err
	}
	
	e.templates[name] = tmpl
	return nil
}

// RenderString renders a template string directly
func (e *Engine) RenderString(templateStr string, data interface{}) (string, error) {
	tmpl, err := template.New("inline").Funcs(e.funcMap).Parse(templateStr)
	if err != nil {
		return "", err
	}
	
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	
	return buf.String(), nil
}

// defaultFuncMap returns default template functions
func defaultFuncMap() template.FuncMap {
	return template.FuncMap{
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": strings.Title,
		"join":  strings.Join,
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"mul": func(a, b int) int {
			return a * b
		},
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"mod": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a % b
		},
		"eq": func(a, b interface{}) bool {
			return a == b
		},
		"ne": func(a, b interface{}) bool {
			return a != b
		},
		"lt": func(a, b int) bool {
			return a < b
		},
		"le": func(a, b int) bool {
			return a <= b
		},
		"gt": func(a, b int) bool {
			return a > b
		},
		"ge": func(a, b int) bool {
			return a >= b
		},
		"default": func(defaultValue, value interface{}) interface{} {
			if value == nil || value == "" {
				return defaultValue
			}
			return value
		},
	}
}
