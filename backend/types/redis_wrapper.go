package types

type ListEntryItem struct {
	Value        any    `json:"v"`
	DisplayValue string `json:"dv,omitempty"`
}

type ListReplaceItem struct {
	Index        int64  `json:"index"`
	Value        any    `json:"v"`
	DisplayValue string `json:"dv,omitempty"`
}

type HashEntryItem struct {
	Key          string `json:"k"`
	Value        any    `json:"v"`
	DisplayValue string `json:"dv,omitempty"`
}

type HashReplaceItem struct {
	Key          any    `json:"k"`
	NewKey       any    `json:"nk"`
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

type ZSetReplaceItem struct {
	Score        float64 `json:"s"`
	Value        string  `json:"v"`
	NewValue     string  `json:"nv"`
	DisplayValue string  `json:"dv,omitempty"`
}

type StreamEntryItem struct {
	ID           string         `json:"id"`
	Value        map[string]any `json:"v"`
	DisplayValue string         `json:"dv,omitempty"`
}
