package main

import (
	"flag"
	"log"
	"os"

	"github.com/abbeydabiri/horizonextractor/extractor"
	"github.com/abbeydabiri/horizonextractor/loader"
	"github.com/abbeydabiri/horizonextractor/transformer"
)

func main() {
	var sFile, sHost, sPort, sUser, sPass, sDB string
	flag.StringVar(&sFile, "file", "", "horizon tool excel file")
	flag.StringVar(&sHost, "host", "", "database host or ip")
	flag.StringVar(&sPort, "port", "", "port number")
	flag.StringVar(&sUser, "user", "", "username")
	flag.StringVar(&sPass, "pass", "", "password")
	flag.StringVar(&sDB, "db", "", "database name")
	flag.Parse()

	sLog := "extractorlog.log"
	fCreate, err := os.Create(sLog)
	fCreate.Close()
	fLog, err := os.OpenFile(sLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(fLog)

	loader.Connect(sHost, sPort, sUser, sPass, sDB)
	xFile, err := extractor.Extract(sFile)
	if err != nil {
		return
	}

	transformer.Transform(xFile)

}
