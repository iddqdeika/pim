package pim

import (
	"fmt"
	"strings"
)

const (
	ArticlePath           = "Article"
	ArticleAttributesPath = "Article/ArticleAttribute"
)

var (
	ArticleAttributesFields = []string{"ArticleAttributeLang.Name", "ArticleAttributeValue.Value(Russian,DEFAULT)"}
	bracketsReplacer        = strings.NewReplacer("[", "", "]", "")
)

type ArticleProvider struct {
	c *Client
}

func (p *ArticleProvider) GetStructureMaps(articleIdentifier string, structureIdentifier string) ([]string, error) {
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

type ArticleStructureMap struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	EntityID int    `json:"entityId"`
}

func (p *ArticleProvider) GetAttributes(articleIdentifier string) ([]ArticleAttribute, error) {
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

type ArticleAttribute struct {
	Name  string
	Value string
}

func (p *ArticleProvider) Update(columns []string, articles ...ArticleUpdate) error {
	url := p.c.baseListUrl() + ArticlePath
	ub, err := newArticleUpdate(columns, articles)
	if err != nil {
		return nil
	}
	res, err := p.c.update(url, ub)
	if err != nil {
		return err
	}
	if res.Counters.Errors != 0 {
		return fmt.Errorf("update complete with %v errors", res.Counters.Errors)
	}
	return nil
}

func (p *ArticleProvider) NewAttrituteValueUpdateOrder(update ArticleAttributeValueUpdate) (*PimUpdateOrder, error) {
	b, err := newArticleAttributeValueUpdate(update)
	if err != nil {
		return nil, err
	}
	return &PimUpdateOrder{
		UrlPath:    ArticleAttributesPath,
		UpdateBody: b,
	}, nil
}

func newArticleAttributeValueUpdate(update ArticleAttributeValueUpdate) (*PimUpdateBody, error) {
	//тело апдейта
	ub := &PimUpdateBody{
		Columns: []PimUpdateColumn{{
			Identifier: "ArticleAttributeValue.Value",
		}},
	}
	//перебираем заданные атрибуты
	for name, value := range update.Attributes {
		//добавляем строчку с объектом значения атрибута
		ub.Rows = append(ub.Rows, PimUpdateRow{
			Object: PimUpdateObject{
				ID: "'" + update.ArticleNo + "'@1"},
			Qualification: map[string]string{
				"name":       name,
				"language":   update.Language,
				"identifier": update.Identifier,
			},
			Values: []string{value},
		})
	}
	return ub, nil
}

func newArticleUpdate(columns []string, articles []ArticleUpdate) (*PimUpdateBody, error) {
	ub := &PimUpdateBody{}
	cm := make(map[string]int)
	for i, column := range columns {
		cm[column] = i
		ub.Columns = append(ub.Columns, PimUpdateColumn{Identifier: column})
	}
	for _, article := range articles {
		row := PimUpdateRow{
			Object: PimUpdateObject{
				ID: "'" + article.ArticleNo + "'@1",
			},
			Values: make([]string, len(cm)),
		}
		if len(article.Fields) != len(cm) {
			return nil, fmt.Errorf("article %v contains %v fields, but there must be %v",
				article.ArticleNo, len(article.Fields), len(cm))
		}
		for field, val := range article.Fields {
			i, ok := cm[field]
			if !ok {
				return nil, fmt.Errorf("article %v contains field %v that isn't declared in fields slice",
					article.ArticleNo, field)
			}
			row.Values[i] = val
		}
		ub.Rows = append(ub.Rows, row)
	}
	return ub, nil
}

func (p *ArticleProvider) NewUpdateFromNo(articleNo string) ArticleUpdate {
	return ArticleUpdate{
		ArticleNo: articleNo,
		Fields:    make(map[string]string),
	}
}

// создать новый объект для обновления значений атрибутов позиции.
// язык и идентификатор значения указывается для всех атрибутов сразу.
func (p *ArticleProvider) NewAttributeValueUpdate(articleNo, language, identifier string) ArticleAttributeValueUpdate {
	return ArticleAttributeValueUpdate{
		ArticleNo:  articleNo,
		Language:   language,
		Identifier: identifier,
		Attributes: make(map[string]string),
	}
}

type ArticleUpdate struct {
	ArticleNo string
	Fields    map[string]string
}

func (a ArticleUpdate) With(field string, value string) ArticleUpdate {
	a.Fields[field] = value
	return a
}

type ArticleAttributeValueUpdate struct {
	ArticleNo  string
	Language   string
	Identifier string
	Attributes map[string]string
}

func (u ArticleAttributeValueUpdate) With(attributeName string, value string) ArticleAttributeValueUpdate {
	u.Attributes[attributeName] = value
	return u
}
