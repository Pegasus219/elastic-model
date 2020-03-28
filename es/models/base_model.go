package models

import (
	"context"
	"encoding/json"
	"errors"
	"gopkg.in/olivere/elastic.v6"
	"reflect"
	"strings"
)

const (
	MAX_SIZE = 10000

	DISTINCT = "distinct"

	AGGS_NAME = "aggs_result"

	SUB_AGGS_NAME = "sub_aggs"
)

type (
	BaseDoc interface {
		getId() string
	}

	BaseModel struct {
		elastic *elastic.Client
		ctx     context.Context
		esIndex string
		esType  string
	}
)

//检查索引是否存在
func isIndexExists(err error) bool {
	if strings.Index(err.Error(), "index_not_found_exception") >= 0 {
		return false
	}
	return true
}

func (m *BaseModel) CheckExists() (bool, error) {
	_, err := m.elastic.Search().Index(m.esIndex).Type(m.esType).Size(0).Do(m.ctx)
	if err != nil {
		exists := isIndexExists(err)
		return exists, err
	}
	return true, nil
}

//创建索引
func (m *BaseModel) CreateIndex(mapping string) error {
	createIndex, err := m.elastic.CreateIndex(m.esIndex).BodyString(mapping).Do(m.ctx)
	if err != nil {
		return err
	}
	if !createIndex.Acknowledged {
		return errors.New("Create Index Error")
	}
	return nil
}

//创建/修改单个文档
func (m *BaseModel) IndexDoc(doc BaseDoc) error {
	id := doc.getId()
	_, err := m.elastic.Index().Index(m.esIndex).Type(m.esType).Id(id).BodyJson(doc).Do(m.ctx)
	return err
}

//批量创建文档
func (m *BaseModel) BatchIndexDoc(docList []BaseDoc) error {
	bulkRequest := m.elastic.Bulk()
	for _, v := range docList {
		id := v.getId()
		req := elastic.NewBulkIndexRequest().Index(m.esIndex).Type(m.esType).Id(id).Doc(v)
		bulkRequest.Add(req)
	}
	_, err := bulkRequest.Do(m.ctx)
	return err
}

//删除文档
func (m *BaseModel) DelDoc(id string) error {
	_, err := m.elastic.Delete().Index(m.esIndex).Type(m.esType).Id(id).Do(m.ctx)
	return err
}

//计数
func (m *BaseModel) doCount(query elastic.Query) (int64, error) {
	searchResult, err := m.elastic.Search().Index(m.esIndex).Type(m.esType).Size(0).Query(query).Do(m.ctx)
	if err != nil {
		return 0, err
	}
	total := searchResult.TotalHits()
	return total, nil
}

//查找
func (m *BaseModel) doSearch(query elastic.Query, limit int, rowStruct interface{}) ([]interface{}, error) {
	searchResult, err := m.elastic.Search().Index(m.esIndex).Type(m.esType).Size(limit).Query(query).Do(m.ctx)
	if err != nil {
		return nil, err
	}
	typ := reflect.TypeOf(rowStruct)
	return searchResult.Each(typ), nil
}

//聚合统计
func (m *BaseModel) doAggregation(query elastic.Query, aggs elastic.Aggregation, result interface{}) error {
	esService := m.elastic.Search().Index(m.esIndex).Type(m.esType).Size(0).Aggregation(AGGS_NAME, aggs)
	if query != nil {
		esService = esService.Query(query)
	}
	searchResult, err := esService.Do(m.ctx)
	if err != nil {
		return err
	}
	bytes, err := searchResult.Aggregations[AGGS_NAME].MarshalJSON()
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, result)
	if err != nil {
		return err
	}
	return nil
}
