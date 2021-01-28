package pim

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"time"
)

const (
	restUrlPath     = "/rest/V1.0/"
	restUrlListPath = "/rest/V1.0/list/"
	retryTimeout    = time.Second
)

var (
	//errors
	TypeCastErr = fmt.Errorf("cant cast value to correct type")
)

type Config struct {
	Host          string
	Login         string
	Password      string
	TimeoutInSecs int
}

func NewClient(config Config) (*Client, error) {
	httpClient := http.Client{
		Transport:     http.DefaultTransport,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * time.Duration(config.TimeoutInSecs),
	}

	c := &Client{
		Config: config,
		client: httpClient,
	}

	// провайдер структурных групп
	sgp := &structureGroupProvider{
		c: c,
	}
	// провайдер позиций
	c.ap = newArticleProvider(c)
	// апдейтер позиций
	c.au = newArticleUpdater(c)
	// провайдер проверок
	dqp := &DataQualityProvider{
		c: c,
	}
	c.sgp = sgp
	c.dqp = dqp
	return c, nil
}

// PIM API client
// provides concrete functional abstractions
type Client struct {
	Config Config
	client http.Client

	sgp StructureGroupProvider
	ap  ArticleProvider
	au  ArticleUpdater
	dqp *DataQualityProvider
}

// позволяет работать со структурными группами
func (c *Client) StructureGroupProvider() StructureGroupProvider {
	return c.sgp
}

// позволяет работать с позициями
func (c *Client) ArticleProvider() ArticleProvider {
	return c.ap
}

func (c *Client) ArticleUpdater() ArticleUpdater {
	return c.au
}

// позволяет работать с правилами качества данных
func (c *Client) DataQualityProvider() *DataQualityProvider {
	return c.dqp
}

func (c *Client) baseUrl() string {
	return "http://" + c.Config.Host + restUrlPath
}

func (c *Client) baseListUrl() string {
	return "http://" + c.Config.Host + restUrlListPath
}

func (c *Client) postJson(url string, data []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.Config.Login, c.Config.Password)
	return c.client.Do(req)
}

func (c *Client) get(url string) (*PimReadResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.Config.Login, c.Config.Password)

	res, err := c.doWithRetries(req, 3)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	r := &PimReadResponse{}
	err = json.Unmarshal(data, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Client) update(url string, ub *PimUpdateBody) (*PimUpdateResponse, error) {
	data, err := json.Marshal(ub)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.Config.Login, c.Config.Password)
	res, err := c.doWithRetries(req, 1)
	if err != nil {
		return nil, err
	}
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	r := &PimUpdateResponse{}
	err = json.Unmarshal(data, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Client) doWithRetries(req *http.Request, tries int) (*http.Response, error) {
	var res *http.Response
	var err error
	for i := 0; i < tries; i++ {
		res, err = c.try(req)
		if err == nil {
			return res, nil
		}
		time.Sleep(retryTimeout)
	}
	return nil, err
}

func (c *Client) try(req *http.Request) (*http.Response, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("status code is %v, body: %v", res.StatusCode, string(body))
	}
	return res, nil
}

func (c *Client) UpdateFromOrder(dto *PimUpdateOrder) error {
	if err := dto.Validate(); err != nil {
		return fmt.Errorf("invalid pim update order given: %v", err)
	}
	url := c.baseListUrl() + dto.UrlPath
	res, err := c.update(url, dto.UpdateBody)
	if err != nil {
		return err
	}
	if res.Counters.Errors != 0 {
		return fmt.Errorf("update complete with %v errors", res.Counters.Errors)
	}
	return nil
}

func (c *Client) DoSearch(s Search) (*PimReadResponse, error) {
	if s == nil {
		return nil, fmt.Errorf("Search is nil")
	}
	if len(s.Query()) == 0 {
		return nil, fmt.Errorf("query is empty")
	}
	url := c.baseListUrl() + s.ReportPath() + "/bySearch?" + "query=" + url2.QueryEscape(s.Query()) + "&fields=" + s.Fields() +
		"&pageSize=-1&cacheId=no-cache"
	return c.get(url)
}

type PimReadResponse struct {
	Rows []PimReadRow `json:"rows"`
}

type PimReadRow struct {
	Object PimReadObject `json:"object"`
	Values []interface{} `json:"values"`
}

type PimReadObject struct {
	ID       string `json:"id"`
	EntityID int    `json:"entityId"`
}

type PimUpdateOrder struct {
	UrlPath    string         `json:"url_path"`
	UpdateBody *PimUpdateBody `json:"update_body"`
}

// проверяет консистентность самое себя
func (o PimUpdateOrder) Validate() error {
	if len(o.UrlPath) == 0 {
		return fmt.Errorf("empty url_path")
	}
	if o.UpdateBody == nil {
		return fmt.Errorf("empty(nil) update body")
	}
	if len(o.UpdateBody.Columns) == 0 {
		return fmt.Errorf("no columns in update body")
	}
	if len(o.UpdateBody.Rows) == 0 {
		return fmt.Errorf("no rows in update body")
	}
	return nil
}

type PimUpdateBody struct {
	Columns []PimUpdateColumn `json:"columns"`
	Rows    []PimUpdateRow    `json:"rows"`
}

type PimUpdateColumn struct {
	Identifier string `json:"identifier"`
}

type PimUpdateRow struct {
	Object        PimUpdateObject   `json:"object"`
	Qualification map[string]string `json:"qualification"`
	Values        []string          `json:"values"`
}

type PimUpdateObject struct {
	ID string `json:"id"`
}

type PimUpdateResponse struct {
	Counters PimUpdateCounters `json:"counters"`
	Entries  []interface{}     `json:"entries"`
	Objects  []interface{}     `json:"objects"`
}

type PimUpdateCounters struct {
	Errors              int `json:"errors"`
	Warnings            int `json:"warnings"`
	CreatedObjects      int `json:"created_objects"`
	UpdatedObjects      int `json:"updated_objects"`
	ObjectsWithErrors   int `json:"objects_with_errors"`
	ObjectsWithWarnings int `json:"objects_with_warnings"`
}
