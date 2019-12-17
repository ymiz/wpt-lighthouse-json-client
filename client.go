package wpt_lighthouse_json_client

import (
	"encoding/json"
	"io/ioutil"
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
	jsonUrl := u.String()
	resp, err := http.Get(jsonUrl)
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
	return c.parseResult(tempResult, jsonUrl)
}

func (c Client) parseResult(jsonResult interface{}, jsonUrl string) (*Result, error) {
	a, ok := jsonResult.(map[string]interface{})
	if !ok {
		return nil, ParseError{Message: "top level parse;"}
	}
	categories, ok := a["categories"]
	if !ok {
		return nil, ParseError{Message: "no categories;"}
	}
	castedCategories, ok := categories.(map[string]interface{})
	if !ok {
		return nil, ParseError{Message: "fail cast categories"}
	}
	performance, ok := castedCategories["performance"]
	if !ok {
		return nil, ParseError{Message: "no performance;"}
	}
	castedPerformance, ok := performance.(map[string]interface{})
	if !ok {
		return nil, ParseError{Message: "fail cast performance;"}
	}
	score, ok := castedPerformance["score"]
	if !ok {
		return nil, ParseError{Message: "no score;"}
	}
	castedScore, ok := score.(float64)
	if !ok {
		return nil, ParseError{Message: "fail cast score;"}
	}
	return &Result{
		Url:        jsonUrl,
		StatusCode: 200,
		Performance: Performance{
			Score: castedScore,
		},
	}, nil
}

type ParseError struct {
	Message string
}

func (p ParseError) Error() string {
	return "lighthouse result parse error; " + p.Message
}

type Params struct {
	TestId string
}
