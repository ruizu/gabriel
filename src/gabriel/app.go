package main

import (
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

var config Config

func init() {
	version := flag.Bool("version", false, "prints version")
	flag.Parse()

	if *version {
		fmt.Println(Version())
		os.Exit(0)
	}

	ok := ReadConfig(&config, "config/gabriel.ini") ||
		ReadConfig(&config, "/etc/gabriel/gabriel.ini") ||
		ReadConfig(&config, "/etc/gabriel.ini")
	if !ok {
		log.Fatal("Could not find configuration file")
	}

	if config.Server.Environment == "production" {
		logger, err := syslog.New(syslog.LOG_NOTICE|syslog.LOG_DAEMON, "gabriel")
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(logger)
		log.SetFlags(0)
	}
}

func main() {
	router := httprouter.New()
	router.GET("/ping", Ping)

	log.Fatal(http.ListenAndServe(config.Server.Host, router))
}

func Ping(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Fprintf(w, "pong")
}
