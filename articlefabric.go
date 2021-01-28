package pim

import "fmt"

var articleFabricInstance ArticleFabric = &articleFabric{}

func GetArticleFabricInstance() ArticleFabric {
	return articleFabricInstance
}

type articleFabric struct {
}

func (p *articleFabric) NewUpdateOrder(columns []string, articles ...ArticleUpdate) (*PimUpdateOrder, error) {
	ub, err := newArticleUpdate(columns, articles)
	if err != nil {
		return nil, err
	}
	return &PimUpdateOrder{
		UrlPath:    ArticlePath,
		UpdateBody: ub,
	}, nil
}

func (p *articleFabric) NewUpdateOrderForAttrituteValue(update ArticleAttributeValueUpdate) (*PimUpdateOrder, error) {
	b, err := newArticleAttributeValueUpdate(update)
	if err != nil {
		return nil, err
	}
	return &PimUpdateOrder{
		UrlPath:    ArticleAttributesPath,
		UpdateBody: b,
	}, nil
}

func (p *articleFabric) NewUpdateFromNo(articleNo string) ArticleUpdate {
	return ArticleUpdate{
		ArticleNo: articleNo,
		Fields:    make(map[string]string),
	}
}

// создать новый объект для обновления значений атрибутов позиции.
// язык и идентификатор значения указывается для всех атрибутов сразу.
func (p *articleFabric) NewUpdateForAttributeValue(articleNo, language, identifier string) ArticleAttributeValueUpdate {
	return ArticleAttributeValueUpdate{
		ArticleNo:  articleNo,
		Language:   language,
		Identifier: identifier,
		Attributes: make(map[string]string),
	}
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
