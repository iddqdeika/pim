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

func (p *ArticleProvider) GetStructureMaps(articleIdentifier string, structureIdentifier string) ([]string, error) {
	url := p.c.baseUrl() + ArticlePath + "/byItems?" +
		"items=" + "'" + articleIdentifier + "'@1" +
		"&fields=ArticleStructureMap.StructureGroup(" + structureIdentifier + ")" +
		"&pageSize=-1"
	res, err := p.c.get(url)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for _, row := range res.Rows {
		for _, v := range row.Values {
			maps, ok := v.([]interface{})
			if !ok {
				continue
			}
			for _, m := range maps {
				mo, ok := m.(map[string]interface{})
				if !ok {
					continue
				}
				li, ok := mo["label"]
				if !ok {
					continue
				}
				label, ok := li.(string)
				if !ok {
					continue
				}
				result = append(result, label)
			}
		}
	}
	return result, nil
}

type ArticleStructureMap struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	EntityID int    `json:"entityId"`
}

//{
//                        "id": "120579@9001",
//                        "label": "170101010101",
//                        "entityId": 3000
//                    }
