package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/lobre/doodle/pkg/models/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	events        *mysql.EventModel
	users         *mysql.UserModel
	templateCache map[string]*template.Template
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	if err := run(infoLog, errorLog); err != nil {
		errorLog.Printf("%s\n", err)
		os.Exit(1)
	}
}

func run(infoLog, errorLog *log.Logger) error {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/doodle?parseTime=true", "MySQL data source name")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "32 bytes secret key for sessions")
	https := flag.Bool("https", false, "Enable HTTPS server")
	flag.Parse()

	db, err := openDB(*dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		return err
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	if *https {
		session.Secure = true
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		events:        &mysql.EventModel{DB: db},
		users:         &mysql.UserModel{DB: db},
		templateCache: templateCache,
	}

	srv := http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if *https {
		srv.TLSConfig = &tls.Config{
			PreferServerCipherSuites: true,
			CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		}

		infoLog.Printf("Starting TLS server on %s", *addr)
		return srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	}

	infoLog.Printf("Starting server on %s", *addr)
	return srv.ListenAndServe()
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
