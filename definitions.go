package pim

type StructureGroupProvider interface {
	GetInfomodelByIdentifier(identifier string, structureID int) (*Infomodel, error)
}

type SearchFabric interface {
	NewSearch(reportPath string) Search
}

type Search interface {
	ReportPath() string
	Query() string
	Fields() string
	WithPredicate(predicate SearchPredicate) Search
	WithOutputField(field string) Search
}

type SearchPredicateFabric interface {
	NewEquals(field string, value string) SearchPredicate
	NewContains(field string, value string) SearchPredicate
	NewContainsIC(field string, value string) SearchPredicate
	Or(p1, p2 SearchPredicate) SearchPredicate
}

type SearchPredicate interface {
	Render() string
}

type Infomodel struct {
	Identifier  string
	ObjectID    string
	StructureID int
	Features    map[string]Feature
}

type Feature struct {
	Name         string
	DataType     string
	PresetValues []string
	Mandatory    bool
	Multivalued  bool
}

type ArticleProvider interface {
	GetStructureMaps(articleIdentifier string, structureIdentifier string) ([]string, error)
	GetAttributes(articleIdentifier string) ([]ArticleAttribute, error)
}

type ArticleUpdater interface {
	DoUpdate(columns []string, articles ...ArticleUpdate) error
}

type ArticleFabric interface {
	NewUpdateOrder(columns []string, articles ...ArticleUpdate) (*PimUpdateOrder, error)
	NewUpdateOrderForAttrituteValue(update ArticleAttributeValueUpdate) (*PimUpdateOrder, error)
	NewUpdateFromNo(articleNo string) ArticleUpdate
	NewUpdateForAttributeValue(articleNo, language, identifier string) ArticleAttributeValueUpdate
}

type ArticleAttribute struct {
	Name  string
	Value string
}
