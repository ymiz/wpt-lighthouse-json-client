package wpt_lighthouse_json_client

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
)

type Client struct {
	BaseUrl string
}

func (c Client) GetLighthouseResult(params Params) (*Result, error) {
	u, err := url.Parse(c.BaseUrl)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join("lighthouse.php")
	q := u.Query()
	q.Set("test", params.TestId)
	q.Set("f", "json")
	u.RawQuery = q.Encode()
	log.Println(u.String())
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result Result
	if resp.StatusCode != 200 {
		result.StatusCode = resp.StatusCode
		return &result, nil
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var tempResult interface{}
	err = json.Unmarshal(bodyBytes, &tempResult)
	if err != nil {
		return nil, err
	}
	score := tempResult.(map[string]interface{})["categories"].(map[string]interface{})["performance"].(map[string]interface{})["score"].(float64)
	result.StatusCode = 200
	result.Performance.Score = score
	return &result, nil
}

type Params struct {
	TestId string
}
