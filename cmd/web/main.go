package main

import (
	"bookingsv2/pkg/config"
	"bookingsv2/pkg/handlers"
	"bookingsv2/pkg/render"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)


const portValue = "localhost:9084"

// create the appconfig var
var app config.AppConfig

var session *scs.SessionManager

// main - main app func
func main() {

	// change to true in production
	app.InProduction = false

	session = scs.New()
	if !app.InProduction {
		session.Cookie.Secure = false
		session.Lifetime = 24 * time.Hour
		session.Cookie.Persist = true
		session.Cookie.SameSite = http.SameSiteLaxMode
	} else {
		session.Cookie.Secure = true
		session.Lifetime = 30 * time.Minute
		session.Cookie.Persist = false
		session.Cookie.SameSite = http.SameSiteDefaultMode
	}

	app.Session = session

	templateCache, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = templateCache

	app.UseCache = false

	repo := handlers.NewRepo(&app)

	handlers.NewHandlers(repo)

	render.NewTemplates(&app)
	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	// if isDevEnv {
	// 	_ = http.ListenAndServe(portValue, nil)
	// } else {
	// 	_ = http.ListenAndServe(":9084", nil)
	// }

	serve := &http.Server {
		Addr: portValue,
		Handler: routes(&app),
	}
	err = serve.ListenAndServe()
	log.Fatal(err)
}