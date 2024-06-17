package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/polyglotdev/vue-api/internal/data"
	"github.com/polyglotdev/vue-api/internal/driver"
)

// config is the type for all application configuration
type config struct {
	port int // what port do we want the web server to listen on
}

// application is the type for all data we want to share with the
// various parts of our application. We will share this information in most
// cases by using this type as the receiver for functions
type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	models   data.Models
}

func main() {
	var cfg config
	cfg.port = 8081

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dsn := os.Getenv("DSN")
	db, err := driver.ConnectPostgres(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		models:   data.New(db.SQL),
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}

// serve starts the web server
func (app *application) serve() error {
	app.infoLog.Printf("API listening on port %d\n", app.config.port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}

	return srv.ListenAndServe()
}
