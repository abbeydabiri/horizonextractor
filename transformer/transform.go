package transformer

import (
	"fmt"
	"log"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/abbeydabiri/horizonextractor/loader"
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

var importID int

//Transform extracts the sheets and their details
func Transform(sFile string, xlFile *excelize.File) {

	lMissing := false
	for iSheet := 1; iSheet <= xlFile.SheetCount; iSheet++ {
		sheetName := xlFile.GetSheetName(iSheet)
		if foundSheets[sheetName] == "missing" {
			foundSheets[sheetName] = "found"
		}
		// println(sheetName)

		// rows, _ := xlFile.GetRows(sheetName)
		// println("rows:", len(rows))

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

	curRow := loader.Import{}
	curRow.Filename = sFile

	//create Record
	sqlInsert, sqlParams := loader.Insert(&curRow, loader.ToMap(&curRow))
	if err := loader.MSSQL.Get(&curRow.ID, sqlInsert, sqlParams...); err != nil {
		log.Println(err.Error())
	}
	importID = curRow.ID

	fmt.Printf("\nExtraction - Process\n")

	//first
	fmt.Printf("\nextracting-transforming-loading: transformSubdomain\n")
	transformSubdomain(xlFile)

	// second
	fmt.Printf("\nextracting-transforming-loading: transformTechnologyref\n")
	transformTechnologyref(xlFile)

	//third
	fmt.Printf("\nextracting-transforming-loading: transformSectorref\n")
	transformSectorref(xlFile)

	//fourth
	fmt.Printf("\nextracting-transforming-loading: transformSubSectorref\n")
	transformSubSectorref(xlFile)

	//fifth
	fmt.Printf("\nextracting-transforming-loading: transformTechrelevancebysector\n")
	transformTechrelevancebysector(xlFile)

	//sixth
	fmt.Printf("\nextracting-transforming-loading: transformTechrelevancebysubsector\n")
	transformTechrelevancebysubsector(xlFile)

	//seventh
	fmt.Printf("\nextracting-transforming-loading: transformTechvulnerabilitybysector\n")
	transformTechvulnerabilitybysector(xlFile)

	//eight
	fmt.Printf("\nextracting-transforming-loading: transformTechvulnerabilitybysubsector\n")
	transformTechvulnerabilitybysubsector(xlFile)

}

func getCellValue(xlFile *excelize.File, sheetName, cellName string) (cellValue string) {
	cellValue, _ = xlFile.GetCellValue(sheetName, cellName)
	return
}
