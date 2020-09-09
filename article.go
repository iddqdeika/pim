package pim

const (
	ArticleWithAttributesPath = "Article/ArticleAttribute"
)

type ArticleProvider struct {
	c *Client
}

func (p *ArticleProvider) MappingAndAttributesByIdentifier(articleIdentifier string, structureIdentifier string) {

}
