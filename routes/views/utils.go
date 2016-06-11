package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)


const (
	resourceDirectory = "routes/views/resources"
	templateDirectory = resourceDirectory + "/templates"
)

func render(w http.ResponseWriter, templateName string, data interface{}) {
	path, _ := filepath.Abs(templateDirectory + "/" + templateName)
	tmpl, _ := template.ParseFiles(path)
	tmpl.Execute(w, data)
}