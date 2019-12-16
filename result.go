package wpt_lighthouse_json_client

type Result struct {
	StatusCode  int
	Performance Performance
}

type Performance struct {
	Score float64
}
