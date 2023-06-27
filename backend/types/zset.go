package types

type ZSetItem struct {
	Value string  `json:"value"`
	Score float64 `json:"score"`
}
