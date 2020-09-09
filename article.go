package pim

import ()

const (
	ArticlePath           = "Article"
	ArticleAttributesPath = "Article/ArticleAttribute"
)

var ()

type ArticleProvider struct {
	c *Client
}

func (p *ArticleProvider) GetStructureMaps(articleIdentifier string, structureIdentifier string) []string {
	url := p.c.baseUrl() + ArticlePath + "/byItems?" +
		"items=" + "'" + articleIdentifier + "'@1" +
		"&fields=ArticleStructureMap.StructureGroup(" + structureIdentifier + ")" +
		"&pageSize=-1"
	res, err := p.c.get(url)
	if err != nil {
		return nil, err
	}
}
