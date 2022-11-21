package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/duypham9895/bookings/pkg/config"
	"github.com/duypham9895/bookings/pkg/handlers"
	"github.com/duypham9895/bookings/pkg/render"

	"github.com/alexedwards/scs/v2"
)

const PORT = ":8080"

var (
	app     config.AppConfig
	session *scs.SessionManager
)

func main() {
	// change this to true when in production
	app.IsUsedCache = false
	app.IsEnvProd = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.IsEnvProd
	app.Session = session

	// Config global variables for Templates
	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache", err.Error())
	}
	app.TemplateCache = templateCache
	render.NewTemplates(&app)

	// Config global variables for Handlers
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	fmt.Println("Starting application on port", PORT)

	server := &http.Server{
		Addr:    PORT,
		Handler: Routes(&app),
	}

	err = server.ListenAndServe()
	log.Fatal(err.Error())
}
