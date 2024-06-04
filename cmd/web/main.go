package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jantoniogonzalez/lets-go/internal/models"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")

	dsn := flag.String("db", "web:lorenzo@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	// Defering it will make the connection close exiting main
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infolog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infolog.Printf("Starting up the server in port %s\n", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}

func openDB(dsn string) (*sql.DB, error) {
	// Doesn't create any connections, only initializes the pool
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// Check for any errors when creating a connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// CONTINUE ON PAGE 121
