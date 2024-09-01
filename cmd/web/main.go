package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"snippetbox.natenine.com/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

// Creating application struct to hold application-wide dependencies for the web application
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// Initializing flag variables
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Define a new command-line flag for the MySQL DSN string
	dsn := flag.String("dsn", "web:Oldsuccess1578@@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	// Initializing logging information
	infoLog := log.New(os.Stdout, "Info\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)

	infoLog.Print("Trying to connect to -> ", *dsn)
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Print("Database Connection Successful")
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal()
	}
	// Initializing a new instance of our application struct, containing the dependencies.
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// Initializing mux to a new server mux

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Logging and starting the server
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
