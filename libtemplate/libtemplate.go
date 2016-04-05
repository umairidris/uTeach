// Package libtemplate provides template related functions.
package libtemplate

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// MarkdownToHTML converts and sanitizes a markdown string into HTML.
func MarkdownToHTML(markdown string) template.HTML {
	unsafe := blackfriday.MarkdownBasic([]byte(markdown))
	safe := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return template.HTML(safe)
}

func FormatAndLocalizeTime(t time.Time) string {
	return t.Local().Format("Jan 2 2006 3:04PM")
}

// http://stackoverflow.com/questions/18276173/calling-a-template-with-several-pipeline-parameters
func Dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}

	dict := map[string]interface{}{}

	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

// Load gets all templates at path into a mapping of the template name to its template object.
// The path should contain a layouts/ subdirectory with all the templates. The path should also contain a includes/
// subdirectory which contains parent and reusable templates, they will be parsed with each template in the layouts/
// directory. See: https://elithrar.github.io/article/approximating-html-template-inheritance/ for implementation
// details.
func Load(path string) (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template)

	funcMap := template.FuncMap{
		"dict":                  Dict,
		"formatAndLocalizeTime": FormatAndLocalizeTime,
		"markdownToHTML":        MarkdownToHTML,
	}

	layouts, err := filepath.Glob(filepath.Join(path, "layouts/*.html"))
	if err != nil {
		return nil, err
	}

	includes, err := filepath.Glob(filepath.Join(path, "includes/*.html"))
	if err != nil {
		return nil, err
	}

	for _, layout := range layouts {
		files := append(includes, layout)
		templates[filepath.Base(layout)] = template.Must(template.New(layout).Funcs(funcMap).ParseFiles(files...))
	}

	return templates, err
}

func Render(w http.ResponseWriter, templates map[string]*template.Template, name string, data map[string]interface{}) error {
	tmpl, ok := templates[name]
	if !ok {
		return errors.New(fmt.Sprintf("The template %s does not exist.", name))
	}

	if data == nil {
		data = map[string]interface{}{}
	}

	// TODO: to speed this up use a buffer pool (https://elithrar.github.io/article/using-buffer-pools-with-go/)
	buf := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return err
	}
	buf.WriteTo(w)
	return nil
}
