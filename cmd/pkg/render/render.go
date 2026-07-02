package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/doha-ms/bnb-booking-app/cmd/pkg/config"
	"github.com/doha-ms/bnb-booking-app/cmd/pkg/models"
)

// tc (Template Cache) is a global locker.
// It stores parsed templates in RAM so we don't read the disk on every request.
// Mechanical Note: map[string] is the filename, *template.Template is the parsed result.
var tc = make(map[string]*template.Template)

// RenderTemplateTest is for development only.
// It bypasses the cache to allow for "Live Reloading" of HTML changes,
// but it is slow because it hits the hard drive every time.
// func RenderTemplateTest(w http.ResponseWriter, tmpl string) {
// 	parsedTemplate, _ := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.tmpl")
// 	err := parsedTemplate.Execute(w, nil)
// 	if err != nil {
// 		fmt.Println("error parsing template", err)
// 		return
// 	}
// }

// RenderTemplate is the production-ready function.
// Flow: Check RAM -> If empty, Parse & Save -> Then Serve.
// func RenderTemplate(w http.ResponseWriter, t string) {
// 	var tmpl *template.Template
// 	var err error

// 	// Check if the template is already in the 'tc' map (the locker)
// 	_, inMap := tc[t]

// 	if !inMap {
// 		log.Println("creating temp and adding to cache")
// 		// Cache Miss: Go build the template from files
// 		err = createTemplateCache(t)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	} else {
// 		// Cache Hit: Pulling directly from RAM (very fast)
// 		log.Println("using cached templates")
// 	}

// 	// Pull the prepared template from our map
// 	tmpl = tc[t]

// 	// Execute "pours" the template into the ResponseWriter (the user's browser)
// 	err = tmpl.Execute(w, nil)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// createTemplateCache does the "heavy lifting" of reading files.
// It combines the specific page (home/about) with the base layout.
// func createTemplateCache(t string) error {
// 	// List all files needed for this specific page
// 	templates := []string{
// 		fmt.Sprintf("./templates/%s", t),
// 		"./templates/base.layout.tmpl",
// 	}

// 	// ParsFiles reads the disk. The '...' unpacks the slice into individual arguments.
// 	tmpl, err := template.ParseFiles(templates...)
// 	if err != nil {
// 		return err // Return error to caller if file path is wrong or HTML is broken
// 	}

// 	// Save the finished result into our global map
// 	tc[t] = tmpl

// 	return nil // nil means "No errors, success!"
// }

var app *config.AppConfig

func NewTemplate(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {

	return td
}
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("couldnt get template from ")
	}
	buf := new(bytes.Buffer)
	td = AddDefaultData(td)
	err := t.Execute(buf, td)
	if err != nil {
		log.Println("Error executing template:", err) // CHECK THIS LOG!
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).ParseFiles(page)
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
