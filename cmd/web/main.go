package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"snippetbox_ilx/internal/models"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP Network Address")
	dsn := flag.String("DSN", "web:ilia2323@/snippetbox?parseTime=true", "database DSN (DRIVER SOURCE NAME)")
	flag.Parse()

	// Creating log files
	logFile, err := os.OpenFile("./log.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer logFile.Close()
	// END of Logging Files

	infoLog := log.New(logFile, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(logFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		snippets: &models.SnippetModel{DB: db},
	}

	srv := http.Server{
		Addr:     *addr,
		Handler:  app.route(),
		ErrorLog: errorLog,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	infoLog.Printf("Starting server on %s", *addr)
	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errorLog.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// gracefully shutdown and ending the log file
	<-sigChan
	infoLog.Print("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		errorLog.Printf("Server Shutdown: %v", err)
	} else {
		infoLog.Print("Server stopped gracefully")
	}
	logFile.WriteString("\n---------------------------------------- END ----------------------------------------\n")
	//END of gracefully shutdown

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
