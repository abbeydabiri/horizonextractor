package transformer

import (
	"fmt"
	"log"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const sheetDisruptivecompanies = "Disruptive companies"
const sheetImpactscores = "Impact scores"
const sheetRelevancescoring = "Relevance scoring"
const sheetSectoraggregate = "Sector aggregate"
const sheetSectorreferences = "Sector references"
const sheetTechnologyreferences = "Technology references"
const sheetTechrelevancebysector = "Tech relevance by sector"
const sheetTechvulnerabilitybysector = "Tech vulnerability by sector"
const sheetVulnerabilityscoring = "Vulnerability scoring"

var foundSheets = map[string]string{
	sheetDisruptivecompanies:       "missing",
	sheetImpactscores:              "missing",
	sheetRelevancescoring:          "missing",
	sheetSectoraggregate:           "missing",
	sheetSectorreferences:          "missing",
	sheetTechnologyreferences:      "missing",
	sheetTechrelevancebysector:     "missing",
	sheetTechvulnerabilitybysector: "missing",
	sheetVulnerabilityscoring:      "missing",
}

//Transform extracts the sheets and their details
func Transform(xlFile *excelize.File) {

	lMissing := false
	for iSheet := 1; iSheet <= xlFile.SheetCount; iSheet++ {
		sheetName := xlFile.GetSheetName(iSheet)
		if foundSheets[sheetName] == "missing" {
			foundSheets[sheetName] = "found"
		}
	}

	for sheet, value := range foundSheets {
		if value == "missing" {
			lMissing = true
			msg := fmt.Sprintf("sheet: %v is missing \n", sheet)
			log.Println(msg)
			fmt.Printf(msg)
		}
	}

	if lMissing {
		return
	}

	//transformation order is important

	//first
	transformSubdomain(xlFile)

	//second
	transformTechnologyref(xlFile)

	//third
	transformSectorref(xlFile)

	//fourt
	transformSubSectorref(xlFile)

}

func getCellValue(xlFile *excelize.File, sheetName, cellName string) (cellValue string) {
	cellValue, _ = xlFile.GetCellValue(sheetName, cellName)
	return
}
