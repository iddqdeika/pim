package pim

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	restUrlPath  = "/rest/V1.0/list/"
	retryTimeout = time.Second
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
	sgp := &StructureGroupProvider{
		c: c,
	}
	ap := &ArticleProvider{
		c: c,
	}
	c.ap = ap
	c.sgp = sgp
	return c, nil
}

type Client struct {
	Config Config
	client http.Client

	sgp *StructureGroupProvider
	ap  *ArticleProvider
}

func (c *Client) StructureGroupProvider() *StructureGroupProvider {
	return c.sgp
}

func (c *Client) ArticleProvider() *ArticleProvider {
	return c.ap
}

func (c *Client) baseUrl() string {
	return "http://" + c.Config.Host + restUrlPath
}

func (c *Client) post(url string, data []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
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
		return nil, fmt.Errorf("status code is %v", res.StatusCode)
	}
	return res, nil
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

type PimUpdateBody struct {
	Columns []PimUpdateColumn `json:"columns"`
	Rows    []PimUpdateRow    `json:"rows"`
}

type PimUpdateColumn struct {
	Identifier string `json:"identifier"`
}

type PimUpdateRow struct {
	Object PimUpdateObject `json:"object"`
	Values []string        `json:"values"`
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

//{
//    "counters": {
//        "errors": 0,
//        "warnings": 0,
//        "createdObjects": 0,
//        "updatedObjects": 1,
//        "objectsWithErrors": 0,
//        "objectsWithWarnings": 0
//    },
//    "entries": [],
//    "objects": [
//        {
//            "row": 0,
//            "object": {
//                "id": "1054540@1",
//                "label": "100023210190",
//                "entityId": 1000
//            },
//            "status": [
//                "UPDATED"
//            ]
//        }
//    ]
//}
