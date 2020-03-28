package structs

type (
	//词频统计结果
	Buckets struct {
		Key      string `json:"key"`
		DocCount int    `json:"doc_count"`
		Distinct struct {
			Value int `json:"value"`
		} `json:"distinct"`
	}
	AggsTerms struct {
		Buckets []*Buckets `json:"buckets"`
	}

	//数值统计结果
	IntBuckets struct {
		Key      int `json:"key"`
		DocCount int `json:"doc_count"`
		Distinct struct {
			Value int `json:"value"`
		} `json:"distinct"`
	}
	AggsIntTerms struct {
		Buckets []*IntBuckets `json:"buckets"`
	}
)
