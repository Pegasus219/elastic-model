package structs

type (
	//统计每个作者每天的文档数
	AuthorDailyCountBuckets struct {
		KeyAsString string        `json:"key_as_string"`
		Key         int64         `json:"key"`
		DocCount    int           `json:"doc_count"`
		SubAggs     *AggsIntTerms `json:"sub_aggs"`
	}
	AggsAuthorDailyCount struct {
		Buckets []*AuthorDailyCountBuckets `json:"buckets"`
	}
)
