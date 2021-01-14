package pim

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type DataQualityProvider struct {
	c *Client
}

func (p *DataQualityProvider) ExecuteRules(order DataQueryOrder) (*DataQualityResult, error) {
	url := p.c.baseUrl() + "manage/dataquality/executions"
	data, err := json.Marshal(order)
	if err != nil {
		return nil, fmt.Errorf("cant marshal DQ rule order: %v", err)
	}
	res, err := p.c.postJson(url, data)
	if err != nil {
		return nil, fmt.Errorf("cant do postJson request to execute DQ rule: %v", err)
	}
	defer res.Body.Close()
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("cant read DQ execution result: %v", err)
	}
	dqres := &DataQualityResult{}
	err = json.Unmarshal(data, dqres)
	if err != nil {
		return nil, fmt.Errorf("cant unmarshal DQ execution result: %v", err)
	}
	return dqres, nil
}

func (p *DataQualityProvider) NewSingleArticleDataQualityOrder(rules []string, articleNo string) (DataQueryOrder, error) {
	if len(rules) < 1 {
		return DataQueryOrder{}, fmt.Errorf("cant create DataQueryOrder: no rules defined")
	}
	return DataQueryOrder{
		RuleNames:        rules,
		EntityIdentifier: "Article",
		ReportQuery: ReportQuery{
			Identifier: "bySearch",
			ParameterList: []ReportParameter{
				{
					Key:   "query",
					Value: "Article.SupplierAID = \"" + articleNo + "\"",
				},
			},
		},
	}, nil
}

// запрос на исполнение DataQuality
type DataQueryOrder struct {
	RuleNames        []string    `json:"rules"`
	EntityIdentifier string      `json:"entityIdentifier"`
	ReportQuery      ReportQuery `json:"reportQuery"`
}

// поисковый запрос.
// определяет набор сущностей, для которых будет исполнено правило
type ReportQuery struct {
	Identifier    string            `json:"identifier"`
	ParameterList []ReportParameter `json:"parameterList"`
}

// параметры для поискового запроса
type ReportParameter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// результат исполнения
type DataQualityResult struct {
	RuleIds                 map[string]int `json:"ruleIds"`
	NumberOfSuccessfulItems int            `json:"numberOfSuccessfulItems"`
	NumberOfFailedItems     int            `json:"numberOfFailedItems"`
}
