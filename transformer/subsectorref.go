package transformer

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/abbeydabiri/horizonextractor/loader"
)

var refSubSectorrefID map[string]int

func transformSubSectorref(xlFile *excelize.File) {
	refSubSectorrefID = make(map[string]int)

	rows, err := xlFile.GetRows(sheetSectorreferences)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err.Error())
	}

	// tableListSubSectorReferences := []loader.SubSectorReferences{}
	for id := range rows {
		cellRow := fmt.Sprintf("%v", id+1)
		if getCellValue(xlFile, sheetSectorreferences, "A"+cellRow) == "" || id == 0 || id == 1 {
			continue
		}

		refKeySector := strings.ToLower(strings.TrimSpace(getCellValue(xlFile, sheetSectorreferences, "A"+cellRow)))
		refKeySubSector := strings.ToLower(strings.TrimSpace(getCellValue(xlFile, sheetSectorreferences, "C"+cellRow)))

		curRow := loader.SubSectorReferences{}
		curRow.Name = getCellValue(xlFile, sheetSectorreferences, "D"+cellRow)
		curRow.SubSectorReferenceNumber = getCellValue(xlFile, sheetSectorreferences, "C"+cellRow)
		curRow.ImpactScore, _ = strconv.Atoi(getCellValue(xlFile, sheetSectorreferences, "E"+cellRow))
		curRow.SectorReferenceID = refSectorrefID[refKeySector]

		//create Record
		refSubSectorrefID[refKeySubSector] = curRow.ID
		// tableListSubSectorReferences = append(tableListSubSectorReferences, curRow)
	}
	// fmt.Printf("tableListSubSectorReferences: %+v", tableListSubSectorReferences[0])

}
