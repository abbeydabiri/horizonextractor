package transformer

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/abbeydabiri/horizonextractor/loader"
)

var refSectorrefID map[string]int

func transformSectorref(xlFile *excelize.File) {
	refSectorrefID = make(map[string]int)

	rows, err := xlFile.GetRows(sheetSectoraggregate)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err.Error())
	}

	// tableListSectorReferences := []loader.SectorReferences{}
	for id := range rows {
		cellRow := fmt.Sprintf("%v", id+1)
		if getCellValue(xlFile, sheetSectoraggregate, "A"+cellRow) == "" || id == 0 || id == 1 {
			continue
		}
		refKeySector := strings.ToLower(strings.TrimSpace(getCellValue(xlFile, sheetSectoraggregate, "A"+cellRow)))

		curRow := loader.SectorReferences{}
		curRow.SectorReferenceNumber = getCellValue(xlFile, sheetSectoraggregate, "A"+cellRow)
		curRow.Name = getCellValue(xlFile, sheetSectoraggregate, "B"+cellRow)
		curRow.AverageSectorImpact, err = strconv.ParseFloat(getCellValue(xlFile, sheetSectoraggregate, "C"+cellRow), 64)

		//create Record
		refSectorrefID[refKeySector] = curRow.ID
		// tableListSectorReferences = append(tableListSectorReferences, curRow)
	}
	// fmt.Printf("tableListSectorReferences: %+v", tableListSectorReferences[0])

}
