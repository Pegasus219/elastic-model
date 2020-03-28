package example

import (
	"context"
	"elastic-model/es/models"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

//测试数据生成
func TestEsBatchIndex(t *testing.T) {
	docs := []models.BaseDoc{
		&models.ExampleDoc{
			Id:       "001",
			SubmitAt: time.Now().AddDate(0, 0, -3),
			Author:   &models.Author{3415, "Tom"},
			Content:  "I use golang for programing",
			Comments: []*models.Comment{
				&models.Comment{78, []string{"A"}},
				&models.Comment{81, []string{"B", "C"}},
			},
		},
		&models.ExampleDoc{
			Id:       "002",
			SubmitAt: time.Now().AddDate(0, 0, -3),
			Author:   &models.Author{1718, "Jack"},
			Content:  "es is useful in this project",
			Comments: []*models.Comment{
				&models.Comment{90, []string{"A", "C"}},
				&models.Comment{88, []string{"B"}},
			},
		},
		&models.ExampleDoc{
			Id:       "003",
			SubmitAt: time.Now().AddDate(0, 0, -2),
			Author:   &models.Author{2341, "Raymond"},
			Content:  "php is no longer the best choice because of golang",
			Comments: []*models.Comment{
				&models.Comment{79, []string{"D"}},
			},
		},
	}
	model := models.NewExampleModel(context.Background())
	err := model.BatchIndexDoc(docs)
	fmt.Println("error=", err)
}

//测试单个文档生成
func TestIndexDoc(t *testing.T) {
	doc := &models.ExampleDoc{
		Id:       "004",
		SubmitAt: time.Now().AddDate(0, 0, -1),
		Author:   &models.Author{2341, "Raymond"},
		Content:  "I cannot use java",
		Comments: []*models.Comment{
			&models.Comment{78, []string{"A", "B"}},
			&models.Comment{82, []string{"B", "C"}},
		},
	}
	model := models.NewExampleModel(context.Background())
	err := model.IndexDoc(doc)
	fmt.Println("error=", err)
}

//测试文档搜索
func TestSearch(t *testing.T) {
	model := models.NewExampleModel(context.Background())
	result, err := model.SearchWith("golang")
	fmt.Println("error=", err)
	strByte, _ := json.Marshal(result)
	fmt.Println(string(strByte))
}

//测试统计，7天内每个作者的文档数量
func TestAggs(t *testing.T) {
	startAt := time.Now().AddDate(0, 0, -7)
	endAt := time.Now()
	model := models.NewExampleModel(context.Background())
	result, err := model.AggsAuthorDailyCount(startAt, endAt)
	fmt.Println("error=", err)
	strByte, _ := json.Marshal(result)
	fmt.Println(string(strByte))
}
