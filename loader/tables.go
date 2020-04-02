package loader

//Tables all tables share this common field for easier processing
type Tables interface {
	Tablename() string
}

var sqlTypes = map[string]string{
	"bool":    "bool",
	"time":    "datetime",
	"string":  "text",
	"int":     "int",
	"uint":    "int",
	"int64":   "int",
	"uint32":  "int",
	"uint64":  "int",
	"float32": "float",
	"float64": "float",
}

//SubDomainReferences table
type SubDomainReferences struct {
	Fields
	Domain string `sql:"index"`
}

//Tablename ...
func (table *SubDomainReferences) Tablename() string {
	return "SubDomainReferences"
}

//TechnologyReferences table
type TechnologyReferences struct {
	Fields
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
func (table *TechnologyReferences) Tablename() string {
	return "TechnologyReferences"
}

//SectorReferences table ...(fetched sector aggregate)
type SectorReferences struct {
	Fields
	SectorReferenceNumber string
	Name                  string `sql:"index"`
	AverageSectorImpact   float64
}

//Tablename ...
func (table *SectorReferences) Tablename() string {
	return "SectorReferences"
}

//SubSectorReferences table
type SubSectorReferences struct {
	Fields
	Name                     string `sql:"index"`
	SubSectorReferenceNumber string `sql:"index"`
	ImpactScore              int
	SectorReferenceID        int
}

//Tablename ...
func (table *SubSectorReferences) Tablename() string {
	return "SubSectorReferences"
}

//TechnologyRelevanceBySectors table
type TechnologyRelevanceBySectors struct {
	Fields
	Score                 int
	SectorReferenceID     int
	TechnologyReferenceID int
}

//Tablename ...
func (table *TechnologyRelevanceBySectors) Tablename() string {
	return "TechnologyRelevanceBySectors"
}

//TechnologyRelevanceBySubSectors table
type TechnologyRelevanceBySubSectors struct {
	Fields
	Score                 int
	TechnologyReferenceID int
	SubSectorReferenceID  int
}

//Tablename ...
func (table *TechnologyRelevanceBySubSectors) Tablename() string {
	return "TechnologyRelevanceBySubSectors"
}

//TechnologyVulnerabilityBySectors table
type TechnologyVulnerabilityBySectors struct {
	Fields
	Score                 int
	TechnologyReferenceID int
	SectorReferenceID     int
}

//Tablename ...
func (table *TechnologyVulnerabilityBySectors) Tablename() string {
	return "TechnologyVulnerabilityBySectors"
}

//TechnologyVulnerabilityBySubSectors table
type TechnologyVulnerabilityBySubSectors struct {
	Fields
	Score                 int
	TechnologyReferenceID int
	SubSectorReferenceID  int
}

//Tablename ...
func (table *TechnologyVulnerabilityBySubSectors) Tablename() string {
	return "TechnologyVulnerabilityBySubSectors"
}

//Readiness table
type Readiness struct {
	Fields
	Score                 int
	Period                string
	SubDomainReferenceID  int
	TechnologyReferenceID int
}

//Tablename ...
func (table *Readiness) Tablename() string {
	return "Readiness"
}

//Import table
type Import struct {
	Fields
	Period, Doctype,
	Filename, Filemeta, Filetype,
	Filepath string
}

//Tablename ...
func (table *Import) Tablename() string {
	return "Import"
}
