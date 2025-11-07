package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP Network Address")

	flag.Parse()

	// Creating log files
	infoFile, err := os.OpenFile("./log.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	errorFile, err := os.OpenFile("./error.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer errorFile.Close()
	defer infoFile.Close()
	// END of Logging Files

	infoLog := log.New(infoFile, "INFO\t", log.Ldate|log.Ltime)

	errorLog := log.New(errorFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	srv := http.Server{
		Addr:     *addr,
		Handler:  app.route(),
		ErrorLog: errorLog,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}
