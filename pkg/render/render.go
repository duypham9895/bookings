package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/duypham9895/bookings/pkg/config"
	"github.com/duypham9895/bookings/pkg/models"
)

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func RenderTemplate(w http.ResponseWriter, t string, templateData *models.TemplateData) {
	var cachedTemplate map[string]*template.Template
	var err error

	if app.IsUsedCache {
		cachedTemplate = app.TemplateCache
		err = nil
	} else {
		cachedTemplate, err = CreateTemplateCache()
		if err != nil {
			log.Fatal("Cannot create template cache", err.Error())
		}
	}

	parsedTemplate, isExistedTemplate := cachedTemplate[t]
	if !isExistedTemplate {
		log.Fatalf("Not found template %s in cached template", t)
	}

	buf := new(bytes.Buffer)

	templateData = AddDefaultData(templateData)
	_ = parsedTemplate.Execute(buf, templateData)

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing template to browser:", err.Error())
	}
}

func AddDefaultData(templateData *models.TemplateData) *models.TemplateData {
	return templateData
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	const PATH_PAGE_FILES = "./templates/*.page.html"
	const PATH_LAYOUT_FILES = "./templates/*.layout.html"

	cachedTemplate := make(map[string]*template.Template)

	// get all files with *.page.html
	pages, err := filepath.Glob(PATH_PAGE_FILES)
	if err != nil {
		return cachedTemplate, err
	}

	layouts, err := filepath.Glob(PATH_LAYOUT_FILES)
	if err != nil {
		return cachedTemplate, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		templateSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return cachedTemplate, err
		}

		if len(layouts) > 0 {
			templateSet, err = templateSet.ParseGlob(PATH_LAYOUT_FILES)
			if err != nil {
				return cachedTemplate, err
			}
		}

		cachedTemplate[name] = templateSet
	}

	return cachedTemplate, nil
}
