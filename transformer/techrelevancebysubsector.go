package transformer

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/abbeydabiri/horizonextractor/loader"
)

func transformTechrelevancebysubsector(xlFile *excelize.File) {

	rows, err := xlFile.GetRows(sheetRelevancescoring)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err.Error())
	}

	var sectorList []string
	sectorCellMap := make(map[string]string)
	// tableListTechnologyRelevanceBySubSectors := []loader.TechnologyRelevanceBySubSectors{}
	for id, cols := range rows {
		cellRow := fmt.Sprintf("%v", id+1)
		if getCellValue(xlFile, sheetRelevancescoring, "A"+cellRow) == "" && id != 4 {
			continue
		}

		if id == 4 {
			iChar := 65
			iPreChar := 0

			for colIndex, colValue := range cols {
				if strings.ToLower(strings.TrimSpace(colValue)) == "" {
					continue
				}

				if colIndex%26 == 0 {
					if iPreChar == 0 {
						iPreChar = 65
					} else {
						iPreChar++
					}
					iChar = 65
				}

				nCurIndex := colIndex % 26

				sectorKey := strings.ToLower(strings.TrimSpace(colValue))
				sectorList = append(sectorList, sectorKey)

				if iPreChar != 0 {
					sectorCellMap[sectorKey] = string(iPreChar) + string(iChar+nCurIndex)
				} else {
					sectorCellMap[sectorKey] = string(iChar + nCurIndex)
				}
			}
		}

		if id <= 5 {
			continue
		}

		refKeyTechnology := strings.ToLower(strings.TrimSpace(getCellValue(xlFile, sheetRelevancescoring, "A"+cellRow)))
		for _, refKeySector := range sectorList {
			curRow := loader.TechnologyRelevanceBySubSectors{}

			sChar := sectorCellMap[refKeySector]
			curRow.Score, _ = strconv.Atoi(getCellValue(xlFile, sheetRelevancescoring, sChar+cellRow))
			curRow.SubSectorReferenceID = refSectorrefID[refKeySector]
			curRow.TechnologyReferenceID = refTechnologyrefID[refKeyTechnology]
			//create Record
			// tableListTechnologyRelevanceBySubSectors = append(tableListTechnologyRelevanceBySubSectors, curRow)
		}
	}
	// fmt.Printf("tableListTechnologyRelevanceBySubSectors: %+v", tableListTechnologyRelevanceBySubSectors)

}
