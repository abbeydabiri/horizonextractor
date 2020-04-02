package transformer

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/abbeydabiri/horizonextractor/loader"
)

var refSubdomainID map[string]int

func transformSubdomain(xlFile *excelize.File) {
	subdomainMap := make(map[string]int)
	refSubdomainID = make(map[string]int)

	rows, err := xlFile.GetRows(sheetTechnologyreferences)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(err.Error())
	}

	colIndex := 2
	for id, cols := range rows {
		if cols[colIndex] == "" || id == 0 {
			continue
		}
		colValue := strings.ToLower(strings.TrimSpace(cols[colIndex]))
		subdomainMap[colValue]++
	}

	var subdomainList []string
	for mapValue := range subdomainMap {
		subdomainList = append(subdomainList, mapValue)
	}
	sort.Strings(subdomainList)

	// tableListSubdomains := []loader.SubDomainReferences{}
	for mapValue := range subdomainMap {
		curRow := loader.SubDomainReferences{}
		curRow.ImportID = importID
		curRow.Domain = strings.Title(mapValue)

		//create Record
		sqlInsert, sqlParams := loader.Insert(&curRow, loader.ToMap(&curRow))
		if err := loader.MSSQL.Get(&curRow.ID, sqlInsert, sqlParams...); err != nil {
			log.Println(sqlInsert)
			log.Println(sqlParams)
			log.Println(err.Error())
		}
		refSubdomainID[mapValue] = curRow.ID

		// tableListSubdomains = append(tableListSubdomains, curRow)
	}
	// fmt.Printf("tableListSubdomains: %+v", tableListSubdomains)
}
