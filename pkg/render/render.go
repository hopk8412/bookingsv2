package render

import (
	"bookingsv2/pkg/config"
	"bookingsv2/pkg/models"
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

// NewTemplates sets teh config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(templateData *models.TemplateData) *models.TemplateData {
	return templateData
}

func RenderTemplate(w http.ResponseWriter, tmpl string, templateData *models.TemplateData) {
	var templateCache map[string]*template.Template
	if app.UseCache {
		//instead of using the below, we want to get template cache from the AppConfig
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}
	// create a template cache - we want this such that we don't have to load it again until the app restarts or we force a restart of the cache
	// can do via an app config
	


	// get the requested template from the cache
	myTemplate, ok := templateCache[tmpl]
	if !ok {
		log.Fatal("Could not get template from cache")
	}

	buffer := new(bytes.Buffer)

	templateData = AddDefaultData(templateData)

	_ = myTemplate.Execute(buffer, templateData)

	// render the template
	_, err := buffer.WriteTo(w)
	if err != nil {
		log.Println(err)
	}


	// loads file from templates directory
	// parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl, "./templates/base.layout.html")
	// err := parsedTemplate.Execute(w, nil)
	// if err != nil {
	// 	fmt.Println("error parsing template:", err)
	// 	return
	// }
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	// returning a map of string, pointer to template.Template
	//myCache := make(map[string]*template.Template)

	//another way to create map - create as an empty map
	myCache := map[string]*template.Template{}

	//get all .page.html files first from templates folder
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}
	// we now have slice of strings containing all files ending .page.html
	//range through these files
	for _, page := range pages {
		//page is the full path for the template name - we want just the template name
		name := filepath.Base(page)
		//try to parse file
		templateSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		hits, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCache, err
		}
		if len(hits) > 0 {
			templateSet, err = templateSet.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return myCache, err
			}
		}
		// this loop will deal with *.page.html, see a layout associated and adds to templateSet, and adds all using ParseGlob function, then adds all to our map
		myCache[name] = templateSet
	}
	return myCache, nil
}


// First approach to template cache - disadvantage in having to specify specific reusable templates (ie base.layout.html) in createTemplateCache function
// Better to have any of these additional templates rendered automatically
// var templateCache = make(map[string]*template.Template)
// func RenderTemplate(w http.ResponseWriter, providedTemplate string) {
// 	//take the parsedTemplate and instead of reading from disk every time, store the parsed template into a map
// 	//check to see if it exists in that map and if not, read from disk, parse, and store in the map
// 	var template *template.Template
// 	var err error

// 	// Check to see if we already have the template in our cache (map)
// 	_, inMap := templateCache[providedTemplate]

// 	if !inMap {
// 		//need to create the template
// 		log.Println("Creating template and adding to cache...")
// 		err = createTemplateCache(providedTemplate)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	} else {
// 		//we have template in cache
// 		log.Println("Using cached template...")
// 	}

// 	template = templateCache[providedTemplate]

// 	err = template.Execute(w, nil)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// func createTemplateCache(templateName string) error {
// 	templates := []string{
// 		fmt.Sprintf("./templates/%s", templateName),
// 		"./templates/base.layout.html",
// 	}
// 	//parse the template - takes each entry from templates slice and puts them in as individual strings
// 	template, err := template.ParseFiles(templates...)

// 	if err != nil {
// 		return err
// 	}

// 	//add template to cache map
// 	templateCache[templateName] = template
// 	return nil
// }