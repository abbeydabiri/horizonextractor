package extractor

import (
	"fmt"
	"log"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

//Extract method works on the excel file
func Extract(sFile string) (xlFile *excelize.File, err error) {

	if sFile == "" {
		sMsg := "Invalid File!!:- file cannot be empty"
		fmt.Println(sMsg)
		log.Fatal(sMsg)
		return
	}

	if !strings.HasSuffix(sFile, "xlsx") {
		sMsg := "Invalid File!!:-  make sure filetype is .xlsx"
		fmt.Println(sMsg)
		log.Fatal(sMsg)
		return
	}

	xlFile, err = excelize.OpenFile(sFile)
	if err != nil {
		sMsg := "File Error! - " + err.Error()
		fmt.Println(sMsg)
		log.Fatal(sMsg)
		return
	}

	return
}
