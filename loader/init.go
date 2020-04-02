package loader

import (
	"time"
)

var allTables = map[string]Tables{
	"SubDomainReferences":                 &SubDomainReferences{},
	"TechnologyReferences":                &TechnologyReferences{},
	"SectorReferences":                    &SectorReferences{},
	"SubSectorReferences":                 &SubSectorReferences{},
	"TechnologyRelevanceBySectors":        &TechnologyRelevanceBySectors{},
	"TechnologyRelevanceBySubSectors":     &TechnologyRelevanceBySubSectors{},
	"TechnologyVulnerabilityBySectors":    &TechnologyVulnerabilityBySectors{},
	"TechnologyVulnerabilityBySubSectors": &TechnologyVulnerabilityBySubSectors{},
	"Readiness":                           &Readiness{},
	"Import":                              &Import{},
}

//Init setup all tables
func Init() {
	for _, table := range allTables {
		Create(table)
	}
}

//Fields ..
type Fields struct {
	ID       int `sql:"pk"`
	ImportID int `sql:"index"`

	Createdate, Updatedate time.Time
}
