package pim

import (
	"encoding/json"
	"fmt"
)

//func newArticleUpdater(c *Client) ArticleUpdater {
//	return &articleUpdater{c: c}
//}

type articleUpdater struct {
	c *Client
}

// столбцы должны содержать идентификаторы полей, указанные в articles.
// если у позиций нет какого-либо из столбцов - вернется соответствующая ошибка
// если у позиций есть поле, которое не перечислено в столбцах - вернется соответствующая ошибка
func (p *articleUpdater) DoUpdate(columns []string, articles ...ArticleUpdate) error {
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
		data, _ := json.Marshal(res)
		return fmt.Errorf("update complete with %v errors, response: %v", res.Counters.Errors, data)
	}
	return nil
}
