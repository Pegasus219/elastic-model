package config

import "os"

var EsUrl = "http://127.0.0.1:9200"

func init() {
	if esUrl := os.Getenv("ES_URL"); esUrl != "" {
		EsUrl = esUrl
	}
}
