package wpt_lighthouse_json_client

type Result struct {
	Url string
	StatusCode  int
	Performance Performance
}

type Performance struct {
	Score float64
}
