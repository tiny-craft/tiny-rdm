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

type ZSetEntryItem struct {
	Score        float64 `json:"s"`
	Value        string  `json:"v"`
	DisplayValue string  `json:"dv,omitempty"`
}

type StreamEntryItem struct {
	ID    string         `json:"id"`
	Value map[string]any `json:"value"`
}
