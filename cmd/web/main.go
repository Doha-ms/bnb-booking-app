package main

import (
	"fmt"
	"log"
	"time"

	"github.com/doha-ms/bnb-booking-app/cmd/pkg/config"
	"github.com/doha-ms/bnb-booking-app/cmd/pkg/handlers"
	"github.com/doha-ms/bnb-booking-app/cmd/pkg/render"

	"net/http"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	app.InProduction = false
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cant create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false
	app.Session = session

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplate(&app)

	// Correct way to serve static files in Go:
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fileServer))
	fmt.Printf("starting application at portt %s", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}
