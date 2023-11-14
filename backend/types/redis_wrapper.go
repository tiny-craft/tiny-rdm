package types

type ListEntryItem struct {
	Value        any    `json:"v"`
	DisplayValue string `json:"dv,omitempty"`
}

type HashEntryItem struct {
	Key          string `json:"k"`
	Value        any    `json:"v"`
	DisplayValue string `json:"dv,omitempty"`
}

type SetEntryItem struct {
	Value        any    `json:"v"`
	DisplayValue string `json:"dv,omitempty"`
}

type ZSetItem struct {
	Value string  `json:"value"`
	Score float64 `json:"score"`
}

type StreamItem struct {
	ID    string         `json:"id"`
	Value map[string]any `json:"value"`
}
