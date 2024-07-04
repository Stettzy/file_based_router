package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

// Router
type Router struct {
	Routes    map[string]string
	Directory string
}

// Instantiate new router
func NewRouter(templatesDir string) *Router {
	return &Router{
		Routes:    make(map[string]string),
		Directory: templatesDir,
	}
}

// Add Route to Router
func (r *Router) AddRoute(path string, template string) {
	r.Routes[path] = template
}

// Render HTML templates
func (r *Router) renderTemplate(w http.ResponseWriter, templateFile string) {
	tmplPath := path.Join(r.Directory, templateFile)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Error finding template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	// Check for exact match first
	if templatePath, ok := r.Routes[path]; ok {
		r.renderTemplate(w, templatePath)
		return
	}

	// If no exact match, try to find the longest matching prefix
	var longestMatch string
	var longestMatchTemplate string

	for routePath, templatePath := range r.Routes {
		if path == routePath {
			longestMatch = routePath
			longestMatchTemplate = templatePath
			break
		}
		if strings.HasPrefix(path, routePath) && len(routePath) > len(longestMatch) {
			longestMatch = routePath
			longestMatchTemplate = templatePath
		}
	}

	if longestMatch != "" && path == longestMatch {
		r.renderTemplate(w, longestMatchTemplate)
		return
	}

	// If no match found, return 404
	http.NotFound(w, req)
}

func main() {
	// Router
	router := NewRouter("./templates")

	// Routes
	router.AddRoute("/", "home.html")
	router.AddRoute("/test", "test.html")
	router.AddRoute("/ok/test", "ok/test.html")

	// Serve
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
