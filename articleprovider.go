package pim

import (
	"fmt"
	"strings"
)

var (
	ArticleAttributesFields = []string{"ArticleAttributeLang.Name", "ArticleAttributeValue.Value(Russian,DEFAULT)"}
	bracketsReplacer        = strings.NewReplacer("[", "", "]", "")
)

func newArticleProvider(c *Client) ArticleProvider {
	return &articleProvider{c: c}
}

type articleProvider struct {
	c *Client
}

// возвращает идентификаторы структурных групп, к которым привязана указанная позиция в указанной структуре
func (p *articleProvider) GetStructureMaps(articleIdentifier string, structureIdentifier string) ([]string, error) {
	url := p.c.baseListUrl() + ArticlePath + "/byItems?" +
		"items=" + "'" + articleIdentifier + "'@1" +
		"&fields=ArticleStructureMap.StructureGroup(" + structureIdentifier + ")->StructureGroup.Identifier" +
		"&pageSize=-1"
	res, err := p.c.get(url)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for _, row := range res.Rows {
		if len(row.Values) != 1 {
			return nil, fmt.Errorf("incorrect value count returned from pim")
		}
		maps, ok := row.Values[0].(string)
		if !ok {
			return nil, TypeCastErr
		}
		maps = bracketsReplacer.Replace(maps)
		result = append(result, maps)
	}
	return result, nil
}

func (p *articleProvider) GetAttributes(articleIdentifier string) ([]ArticleAttribute, error) {
	url := p.c.baseListUrl() + ArticleAttributesPath + "/byItems?" +
		"items=" + "'" + articleIdentifier + "'@1" +
		"&fields=" + strings.Join(ArticleAttributesFields, ",") +
		"&pageSize=-1"
	res, err := p.c.get(url)
	if err != nil {
		return nil, err
	}
	result := make([]ArticleAttribute, 0)
	for _, row := range res.Rows {
		if len(row.Values) != len(ArticleAttributesFields) {
			return nil, fmt.Errorf("cant parse attributes, wrong num of values in a row")
		}
		name, ok := row.Values[0].(string)
		if !ok {
			return nil, TypeCastErr
		}
		value, ok := row.Values[1].(string)
		if !ok {
			return nil, TypeCastErr
		}
		result = append(result, ArticleAttribute{
			Name:  name,
			Value: value,
		})
	}
	return result, nil
}
