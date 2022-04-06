package pim

type StructureGroupProvider interface {
	GetInfomodelByIdentifier(identifier string, structureID int) (*Infomodel, error)
}

type SearchFactory interface {
	NewSearch(reportPath string) Search
	NewStructureGroupSearch(structureIdentifier string) Search
}

// Search defines search params and output fields
type Search interface {
	ReportPath() string
	Query() string
	Fields() string
	Params() map[string]string
	WithPredicate(predicate SearchPredicate) Search
	WithOutputField(field string) Search
	WithAdditionalParam(name, value string) Search
}

type SearchPredicateFactory interface {
	Equals(field string, value string) SearchPredicate
	EqualsIC(field string, value string) SearchPredicate
	NotEquals(field string, value string) SearchPredicate
	NotEqualsIC(field string, value string) SearchPredicate
	IsEmpty(field string) SearchPredicate
	NotIsEmpty(field string) SearchPredicate
	Contains(field string, value string) SearchPredicate
	NotContains(field string, value string) SearchPredicate
	ContainsIC(field string, value string) SearchPredicate
	NotContainsIC(field string, value string) SearchPredicate
	Or(p1, p2 SearchPredicate) SearchPredicate
	And(p1, p2 SearchPredicate) SearchPredicate
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
	CheckArticleExistence(articleIdentifier string) (exists bool, err error)
	SetArticleMediaAssets(articleIdentifier string, assets []*Asset) error
	DeleteArticleMediaAssets(articleIdentifier string, assets []*Asset) error
}

type ArticleAttribute struct {
	Name  string
	Value string
}

type AssetProvider interface {
	UploadFile(*File) (*UploadedFile, error)
	RegisterAsset(*UploadedFile) (*Asset, error)
}

// отключено для последующего удаления
//type ArticleUpdater interface {
//	DoUpdate(columns []string, articles ...ArticleUpdate) error
//}

type ArticleUpdateFactory interface {
	NewUpdateOrder(columns []string, articles ...ArticleUpdate) (*PimUpdateOrder, error)
	NewUpdateOrderForAttrituteValue(update ArticleAttributeValueUpdate) (*PimUpdateOrder, error)
	NewUpdateFromNo(articleNo string) ArticleUpdate
	NewUpdateForAttributeValue(articleNo, language, identifier string) ArticleAttributeValueUpdate
	NewDeleteMediaAssetOrder(article ArticleMediaAssetDelete) (*PimDeleteOrder, error)
	NewDeleteFromNo(articleNo string, mediaTypes []string) ArticleMediaAssetDelete
}

type AssetFactory interface {
	NewFile(name string, filetype string, data []byte) *File
}

type StructureGroupUpdateFactory interface {
	NewUpdateOrder(columns []string, groups ...StructureGroupUpdate) (*PimUpdateOrder, error)
}
