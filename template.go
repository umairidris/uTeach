package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

var templates = make(map[string]*template.Template)

func LoadTemplates() {
	baseTemplate := template.Must(template.ParseFiles("tmpl/base.html"))

	layoutFiles, _ := filepath.Glob("tmpl/layout/*.html")
	for _, layoutFile := range layoutFiles {
		baseTemplateCopy, err := baseTemplate.Clone()
		if err != nil {
			panic(err)
		}
		templates[filepath.Base(layoutFile)] = template.Must(baseTemplateCopy.ParseFiles(layoutFile))
	}
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	tmpl, ok := templates[name]
	if !ok {
		http.Error(w, fmt.Sprintf("The template %s does not exist.", tmpl), http.StatusInternalServerError)
		return
	}

	user, ok := GetSessionUser(r)
	if !ok {
		// if failed to get user, make sure user is nil so templates don't render a user
		user = nil
	}

	data["User"] = user

	err := tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
}
