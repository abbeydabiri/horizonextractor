package loader

//Tables all tables share this common field for easier processing
type Tables interface {
	Tablename() string
}

var sqlTypes = map[string]string{
	"bool":    "bool",
	"time":    "timestamp",
	"string":  "text",
	"int":     "int",
	"uint":    "int",
	"int64":   "int8",
	"uint32":  "int8",
	"uint64":  "int8",
	"float32": "float8",
	"float64": "float8",
}

//SubDomainReferences table
type SubDomainReferences struct {
	ID     int    `sql:"pk"`
	Domain string `sql:"index"`
}

//Tablename ...
func (table *SubDomainReferences) Tablename() {}

//TechnologyReferences table
type TechnologyReferences struct {
	ID                        int    `sql:"pk"`
	TechnologyReferenceNumber string `sql:"index"`
	DisruptorName             string `sql:"index"`
	VendorPresence            int
	MarketPresence            int
	TopicMaturity             int
	SubDomainReferenceID      int
	AdoptionCurveTemplate     string `sql:"index"`
	AdoptionHorizon           string `sql:"index"`
	ImpactAcrossSectors       int
	MidpointOfAdoptionHorizon string `sql:"index"`
	Vulnerability             int
	CategoryID                int
}

//Tablename ...
func (table *TechnologyReferences) Tablename() {}

//SectorReferences table ...(fetched sector aggregate)
type SectorReferences struct {
	ID                    int `sql:"pk"`
	SectorReferenceNumber string
	Name                  string `sql:"index"`
	AverageSectorImpact   float64
}

//Tablename ...
func (table *SectorReferences) Tablename() {}

//SubSectorReferences table
type SubSectorReferences struct {
	ID                       int    `sql:"pk"`
	Name                     string `sql:"index"`
	SubSectorReferenceNumber string `sql:"index"`
	ImpactScore              int
	SectorReferenceID        int
}

//Tablename ...
func (table *SubSectorReferences) Tablename() {}

//TechnologyRelevanceBySectors table
type TechnologyRelevanceBySectors struct {
	ID                    int `sql:"pk"`
	Score                 int
	SectorReferenceID     int
	TechnologyReferenceID int
}

//Tablename ...
func (table *TechnologyRelevanceBySectors) Tablename() {}

//TechnologyRelevanceBySubSectors table
type TechnologyRelevanceBySubSectors struct {
	ID                    int `sql:"pk"`
	Score                 int
	TechnologyReferenceID int
	SubSectorReferenceID  int
}

//Tablename ...
func (table *TechnologyRelevanceBySubSectors) Tablename() {}

//TechnologyVulnerabilityBySectors table
type TechnologyVulnerabilityBySectors struct {
	ID                    int `sql:"pk"`
	Score                 int
	TechnologyReferenceID int
	SectorReferenceID     int
}

//Tablename ...
func (table *TechnologyVulnerabilityBySectors) Tablename() {}

//TechnologyVulnerabilityBySubSectors table
type TechnologyVulnerabilityBySubSectors struct {
	ID                    int
	Score                 int
	TechnologyReferenceID int
	SubSectorReferenceID  int
}

//Tablename ...
func (table *TechnologyVulnerabilityBySubSectors) Tablename() {}

//Readiness table
type Readiness struct {
	ID     int
	Period string
	Score  int
}

//Tablename ...
func (table *Readiness) Tablename() {}
