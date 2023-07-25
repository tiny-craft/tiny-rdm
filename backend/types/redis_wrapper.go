package types

type ZSetItem struct {
	Value string  `json:"value"`
	Score float64 `json:"score"`
}

type StreamItem struct {
	ID    string         `json:"id"`
	Value map[string]any `json:"value"`
}
