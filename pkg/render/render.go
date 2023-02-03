package render

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/Jake-Sheehan/bookings/pkg/config"
	"github.com/Jake-Sheehan/bookings/pkg/models"
)

var templateCache = make(map[string]*template.Template)

// NewTemplates sets the config for the template package
var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(templateData *models.TemplateData) *models.TemplateData {
	return templateData
}

func RenderTemplate(w http.ResponseWriter, tmpl string, templateData *models.TemplateData) {
	var err error
	var templateCache map[string]*template.Template
	// get template cache from app config
	if app.UseCache {
		templateCache = app.TemplateCache
	} else {
		templateCache, err = CreateTemplateCache()
		if err != nil {
			log.Fatal(err)
		}
	}

	// get requested template from cache
	template, ok := templateCache[tmpl]
	if !ok {
		log.Fatal(err)
	}

	// write to buffer first to check for err
	buf := new(bytes.Buffer)

	templateData = AddDefaultData(templateData)

	err = template.Execute(buf, templateData)
	if err != nil {
		log.Println(err)
	}

	// write buffer contents to resonse writer
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	templateCache := map[string]*template.Template{}

	// get all pages files
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return templateCache, err
	}

	// range over pages
	for _, page := range pages {
		// get the file name from the full path
		fileName := filepath.Base(page)
		// create new template object and parse the page
		templateSet, err := template.New(fileName).ParseFiles(page)
		if err != nil {
			return templateCache, err
		}

		// get a list of layout files
		layouts, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return templateCache, err
		}

		// if there are layout files, parse those and add to template
		if len(layouts) > 0 {
			templateSet, err = templateSet.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return templateCache, err
			}
		}
		templateCache[fileName] = templateSet
	}
	return templateCache, nil
}
