package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Jake-Sheehan/bookings/pkg/config"
	"github.com/Jake-Sheehan/bookings/pkg/handlers"
	"github.com/Jake-Sheehan/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const port string = ":8080"

var app config.AppConfig

var session *scs.SessionManager

func main() {
	var err error

	// change this to true in production
	app.InProduction = false

	// session settings
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	// create the template cache
	app.TemplateCache, err = render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	// app settings
	app.UseCache = true

	// creates a repo in handlers package
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	// sets new templates in render package
	render.NewTemplates(&app)

	// runs the server
	fmt.Println("server running on port", port)
	srv := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
