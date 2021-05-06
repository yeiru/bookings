package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/yeiru/bookings/pkg/config"
	"github.com/yeiru/bookings/pkg/models"
)

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, tmplName string, tData *models.TemplateData) {
	var templateCache map[string]*template.Template
	if app.UseCache {
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	tc, ok := templateCache[tmplName]
	if !ok {
		log.Fatal("File not found")
	}
	// parseTemplate, _ := template.ParseFiles("./templates/" + tmplName)
	// err = parseTemplate.Execute(w, nil)

	buffer := new(bytes.Buffer)

	tData = AddDefaultData(tData)

	_ = tc.Execute(buffer, tData)

	_, err := buffer.WriteTo(w)
	if err != nil {
		fmt.Println("Error parsing template", err)
	}
}

var functions = template.FuncMap{}

func CreateTemplateCache() (map[string]*template.Template, error) {
	/*parseTemplate, _ := template.ParseFiles("./templates/" + tmplName)
	err := parseTemplate.Execute(w, nil)
	if err != nil {
		fmt.Println("Error parsing template", err)
		return
	}*/

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")

	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}
