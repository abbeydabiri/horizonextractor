package transformer

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/abbeydabiri/horizonextractor/loader"
)

var refTechnologyrefID map[string]int

func transformTechnologyref(xlFile *excelize.File) {
	refTechnologyrefID = make(map[string]int)

	rows, err := xlFile.GetRows(sheetTechnologyreferences)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err.Error())
	}

	// tableListTechnologyReferences := []loader.TechnologyReferences{}
	for id := range rows {
		cellRow := fmt.Sprintf("%v", id+1)
		if getCellValue(xlFile, sheetTechnologyreferences, "A"+cellRow) == "" || id == 0 {
			continue
		}

		refKeyTech := strings.ToLower(strings.TrimSpace(getCellValue(xlFile, sheetTechnologyreferences, "A"+cellRow)))
		refKeySubdomain := strings.ToLower(strings.TrimSpace(getCellValue(xlFile, sheetTechnologyreferences, "C"+cellRow)))

		curRow := loader.TechnologyReferences{}
		curRow.ImportID = importID
		curRow.TechnologyReferenceNumber = getCellValue(xlFile, sheetTechnologyreferences, "A"+cellRow)
		curRow.DisruptorName = getCellValue(xlFile, sheetTechnologyreferences, "B"+cellRow)

		curRow.VendorPresence, _ = strconv.Atoi(getCellValue(xlFile, sheetTechnologyreferences, "P"+cellRow))
		curRow.MarketPresence, _ = strconv.Atoi(getCellValue(xlFile, sheetTechnologyreferences, "Q"+cellRow))
		curRow.TopicMaturity, _ = strconv.Atoi(getCellValue(xlFile, sheetTechnologyreferences, "R"+cellRow))
		curRow.SubDomainReferenceID = refSubdomainID[refKeySubdomain]
		curRow.AdoptionCurveTemplate = getCellValue(xlFile, sheetTechnologyreferences, "J"+cellRow)
		curRow.AdoptionHorizon = getCellValue(xlFile, sheetTechnologyreferences, "H"+cellRow)
		curRow.ImpactAcrossSectors, _ = strconv.Atoi(getCellValue(xlFile, sheetTechnologyreferences, "K"+cellRow))
		curRow.MidpointOfAdoptionHorizon = getCellValue(xlFile, sheetTechnologyreferences, "I"+cellRow)
		curRow.Vulnerability, _ = strconv.Atoi(getCellValue(xlFile, sheetTechnologyreferences, "L"+cellRow))

		//create Record
		sqlInsert, sqlParams := loader.Insert(&curRow, loader.ToMap(&curRow))
		if err := loader.MSSQL.Get(&curRow.ID, sqlInsert, sqlParams...); err != nil {
			log.Println(err.Error())
		}
		refTechnologyrefID[refKeyTech] = curRow.ID

		// tableListTechnologyReferences = append(tableListTechnologyReferences, curRow)

		//After creating technology reference.

		for readinessCounter := 1; readinessCounter <= 4; readinessCounter++ {
			colChar := string(67 + readinessCounter) //char value of 67 is C
			curRowReadiness := loader.Readiness{}
			curRowReadiness.ImportID = importID
			curRowReadiness.TechnologyReferenceID = refTechnologyrefID[refKeyTech]
			curRowReadiness.SubDomainReferenceID = refSubdomainID[refKeySubdomain]
			curRowReadiness.Period = getCellValue(xlFile, sheetTechnologyreferences, colChar+"1")
			curRowReadiness.Score, _ = strconv.Atoi(getCellValue(xlFile, sheetTechnologyreferences, colChar+cellRow))
			sqlInsert, sqlParams := loader.Insert(&curRowReadiness, loader.ToMap(&curRowReadiness))
			if err := loader.MSSQL.Get(&curRowReadiness.ID, sqlInsert, sqlParams...); err != nil {
				log.Println(err.Error())
			}
		}
		//After creating technology reference.

	}
	// fmt.Printf("tableListTechnologyReferences: %+v", tableListTechnologyReferences[0])

}
