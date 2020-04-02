package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/abbeydabiri/horizonextractor/extractor"
	"github.com/abbeydabiri/horizonextractor/loader"
	"github.com/abbeydabiri/horizonextractor/transformer"
)

func main() {
	var sFile, sServer, sDB, sUser, sPass string
	flag.StringVar(&sFile, "file", "", "horizon tool excel file")
	flag.StringVar(&sServer, "server", "", "database server")
	flag.StringVar(&sDB, "db", "", "database name")
	flag.StringVar(&sUser, "user", "", "username")
	flag.StringVar(&sPass, "pass", "", "password")
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

	go func() {
		ticker := time.Tick(time.Second)
		for i := 1; i != 0; i++ {
			<-ticker
			fmt.Printf("\r %d seconds elapsed", i)
		}
	}()

	loader.Connect(sServer, sDB, sUser, sPass)
	if !loader.Connected {
		log.Fatal("Unable to Connect to DB")
		fmt.Println("Unable to Connect to DB")
		return
	}
	loader.Init()

	xFile, err := extractor.Extract(sFile)
	if err != nil {
		return
	}

	transformer.Transform(sFile, xFile)
	sMsg := "\nExtract, Transform and Load completed - pls check records "
	fmt.Println(sMsg)
	log.Printf(sMsg)
}
