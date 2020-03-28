package es

import (
	"elastic-model/config"
	"gopkg.in/olivere/elastic.v6"
)

var client *elastic.Client

//获取ES Client实例
func GetElastic() *elastic.Client {
	if client != nil {
		return client
	}

	var err error
	client, err = elastic.NewClient(
		elastic.SetURL(config.EsUrl),
		elastic.SetSniff(false),
	)
	if err != nil {
		panic(err)
	}
	return client
}
