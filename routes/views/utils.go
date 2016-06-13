package views

import (
	"github.com/crob1140/codewiz/log"
	"github.com/crob1140/codewiz/models/users"
	"github.com/gorilla/sessions"
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
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		log.Error("Failed to parse template", log.Fields{"template" : templateName, "error" : err})
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Error("Failed to execute template", log.Fields{"template" : templateName, "data" : data, "error" : err})
	}
}

func getUserForSession(session *sessions.Session, userDao *users.Dao) (*users.User, error) {
	userID := session.Values["userID"]
	if userID == nil {
		return nil, nil
	}

	return userDao.GetByID(userID.(uint64))
}