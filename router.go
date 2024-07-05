package router

import (
	"html/template"
	"net/http"
)

// Router
type Router struct {
	Routes    map[string]string
	Directory string
	Renderer  func(w http.ResponseWriter, templateFile string) error
}

// Instantiate new router
func NewRouter(templatesDir string) *Router {
	r := &Router{
		Routes:    make(map[string]string),
		Directory: templatesDir,
	}

	r.Renderer = r.defaultTemplateRenderer

	return r
}

func (r *Router) defaultTemplateRenderer(w http.ResponseWriter, templatePath string) error {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "No template found", http.StatusNotFound)
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to open template", http.StatusInternalServerError)
		return err
	}

	return nil
}

// Add Route to Router
func (r *Router) AddRoute(path string, template string) {
	r.Routes[path] = template
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	// Check for exact match
	if templatePath, ok := r.Routes[path]; ok {
		err := r.Renderer(w, templatePath+".html")
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
		return
	}

	// If no match found, return 404
	http.NotFound(w, req)
}
