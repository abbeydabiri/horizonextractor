package transformer

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/abbeydabiri/horizonextractor/loader"
)

func transformTechrelevancebysector(xlFile *excelize.File) {

	rows, err := xlFile.GetRows(sheetTechrelevancebysector)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err.Error())
	}

	var sectorList []string
	sectorCellMap := make(map[string]string)
	// tableListTechnologyRelevanceBySectors := []loader.TechnologyRelevanceBySectors{}
	for id := range rows {
		cellRow := fmt.Sprintf("%v", id+1)
		if getCellValue(xlFile, sheetTechrelevancebysector, "A"+cellRow) == "" && id != 4 {
			continue
		}

		if id == 4 {
			for iChar := 65; iChar < 91; iChar++ {
				colValue := getCellValue(xlFile, sheetTechrelevancebysector, string(iChar)+cellRow)
				if strings.ToLower(strings.TrimSpace(colValue)) == "" {
					continue
				}
				sectorKey := strings.ToLower(strings.TrimSpace(colValue))
				sectorList = append(sectorList, sectorKey)
				sectorCellMap[sectorKey] = string(iChar)
			}
		}

		if id <= 5 {
			continue
		}

		refKeyTechnology := strings.ToLower(strings.TrimSpace(getCellValue(xlFile, sheetTechrelevancebysector, "A"+cellRow)))
		for _, refKeySector := range sectorList {
			curRow := loader.TechnologyRelevanceBySectors{}
			curRow.ImportID = importID
			sChar := sectorCellMap[refKeySector]
			curRow.Score, _ = strconv.Atoi(getCellValue(xlFile, sheetTechrelevancebysector, sChar+cellRow))
			curRow.SectorReferenceID = refSectorrefID[refKeySector]
			curRow.TechnologyReferenceID = refTechnologyrefID[refKeyTechnology]

			//create Record
			sqlInsert, sqlParams := loader.Insert(&curRow, loader.ToMap(&curRow))
			if err := loader.MSSQL.Get(&curRow.ID, sqlInsert, sqlParams...); err != nil {
				log.Println(err.Error())
			}

			// tableListTechnologyRelevanceBySectors = append(tableListTechnologyRelevanceBySectors, curRow)
		}
	}
	// fmt.Printf("tableListTechnologyRelevanceBySectors: %+v", tableListTechnologyRelevanceBySectors)

}
