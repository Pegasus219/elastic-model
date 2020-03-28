package models

import (
	"context"
	"elastic-model/es"
	"elastic-model/es/structs"
	"gopkg.in/olivere/elastic.v6"
	"time"
)

const EXAMPLE_MAPPING = `
{
	"mappings": {
		"example": {
			"properties": {
				"id": {
					"type": "keyword"
				},
				"submitAt": {
					"type": "date"
				},
				"author": {
          			"properties": {
						"id": {
							"type": "integer"
						},
						"name": {
							"type": "keyword"
						}
					}
        		},
				"content": {
					"type": "text",
					"analyzer": "standard"
				},
				"comments": {
					"type":"nested",
					"properties":{  
						"score":{ 
							"type": "byte"
						},
						"keywords":{ 
							"type": "keyword"
						}
					}
				}
			}
		}
	}
}
`

type (
	Author struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	Comment struct {
		Score    byte     `json:"score"`
		Keywords []string `json:"keywords"`
	}

	ExampleDoc struct {
		Id       string     `json:"id"`
		SubmitAt time.Time  `json:"submitAt"`
		Author   *Author    `json:"author"`
		Content  string     `json:"content"`
		Comments []*Comment `json:"comments"`
	}

	ExampleModel struct {
		BaseModel
	}
)

//实现BaseDoc接口
func (doc *ExampleDoc) getId() string {
	return doc.Id
}

//检查并创建索引
func init() {
	m := NewExampleModel(context.Background())
	exist, _ := m.CheckExists()
	if !exist {
		err := m.CreateIndex(EXAMPLE_MAPPING)
		if err != nil {
			panic(err)
		}
	}
}

func NewExampleModel(ctx context.Context) *ExampleModel {
	m := &ExampleModel{}
	m.ctx = ctx
	m.elastic = es.GetElastic()
	m.esIndex = "example"
	m.esType = "example"
	return m
}

//查找包含指定内容的文档
func (m *ExampleModel) SearchWith(text string) ([]*ExampleDoc, error) {
	contentMatch := elastic.NewMatchQuery("content", text)
	query := elastic.NewBoolQuery().Filter(contentMatch)
	ret, err := m.doSearch(query, 10, new(ExampleDoc))
	if err != nil {
		return nil, err
	}
	var result []*ExampleDoc
	for _, v := range ret {
		result = append(result, v.(*ExampleDoc))
	}
	return result, nil
}

//统计每个作者每天的文档数
func (m *ExampleModel) AggsAuthorDailyCount(startAt, endAt time.Time) ([]*structs.AuthorDailyCountBuckets, error) {
	//query
	timeRange := elastic.NewRangeQuery("submitAt").From(startAt).To(endAt)
	query := elastic.NewBoolQuery().Filter(timeRange)
	//aggs
	authorAggs := elastic.NewTermsAggregation().Field("author.id").Size(MAX_SIZE)
	boundsMin := startAt.Unix() * 1000
	boundsMax := endAt.Unix() * 1000
	aggs := elastic.NewDateHistogramAggregation().Field("submitAt").
		Format("yyyy-MM-dd").TimeZone("+08:00").Interval("day").
		MinDocCount(0).ExtendedBounds(boundsMin, boundsMax).
		SubAggregation(SUB_AGGS_NAME, authorAggs)
	//输出聚合统计结果
	var result structs.AggsAuthorDailyCount
	err := m.doAggregation(query, aggs, &result)
	if err != nil {
		return nil, err
	}
	return result.Buckets, nil
}
