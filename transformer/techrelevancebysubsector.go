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

	var subSectorList []string
	subSectorCellMap := make(map[string]string)
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

				subSectorKey := strings.ToLower(strings.TrimSpace(colValue))
				subSectorList = append(subSectorList, subSectorKey)

				if iPreChar != 0 {
					subSectorCellMap[subSectorKey] = string(iPreChar) + string(iChar+nCurIndex)
				} else {
					subSectorCellMap[subSectorKey] = string(iChar + nCurIndex)
				}
			}
		}

		if id <= 5 {
			continue
		}

		refKeyTechnology := strings.ToLower(strings.TrimSpace(getCellValue(xlFile, sheetRelevancescoring, "A"+cellRow)))
		for _, refKeySubSector := range subSectorList {
			curRow := loader.TechnologyRelevanceBySubSectors{}
			curRow.ImportID = importID
			sChar := subSectorCellMap[refKeySubSector]
			curRow.Score, _ = strconv.Atoi(getCellValue(xlFile, sheetRelevancescoring, sChar+cellRow))
			curRow.SubSectorReferenceID = refSubSectorrefID[refKeySubSector]
			curRow.TechnologyReferenceID = refTechnologyrefID[refKeyTechnology]
			//create Record
			sqlInsert, sqlParams := loader.Insert(&curRow, loader.ToMap(&curRow))
			if err := loader.MSSQL.Get(&curRow.ID, sqlInsert, sqlParams...); err != nil {
				log.Println(err.Error())
			}

			// tableListTechnologyRelevanceBySubSectors = append(tableListTechnologyRelevanceBySubSectors, curRow)
		}
	}
	// fmt.Printf("tableListTechnologyRelevanceBySubSectors: %+v", tableListTechnologyRelevanceBySubSectors)

}
