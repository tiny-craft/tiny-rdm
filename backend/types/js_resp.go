package types

type JSResp struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Data    any    `json:"data,omitempty"`
}

type KeySummaryParam struct {
	Server string `json:"server"`
	DB     int    `json:"db"`
	Key    any    `json:"key"`
}

type KeySummary struct {
	Type   string `json:"type"`
	TTL    int64  `json:"ttl"`
	Size   int64  `json:"size"`
	Length int64  `json:"length"`
}

type KeyDetailParam struct {
	Server       string `json:"server"`
	DB           int    `json:"db"`
	Key          any    `json:"key"`
	ViewAs       string `json:"viewAs,omitempty"`
	DecodeType   string `json:"decodeType,omitempty"`
	MatchPattern string `json:"matchPattern,omitempty"`
	Reset        bool   `json:"reset"`
	Full         bool   `json:"full"`
}

type KeyDetail struct {
	Value      any    `json:"value"`
	Length     int64  `json:"length,omitempty"`
	ViewAs     string `json:"viewAs,omitempty"`
	DecodeType string `json:"decodeType,omitempty"`
	End        bool   `json:"end"`
}

type SetKeyParam struct {
	Server  string `json:"server"`
	DB      int    `json:"db"`
	Key     any    `json:"key"`
	KeyType string `json:"keyType"`
	Value   any    `json:"value"`
	TTL     int64  `json:"ttl"`
	Format  string `json:"format,omitempty"`
	Decode  string `json:"decode,omitempty"`
}

type SetHashParam struct {
	Server   string `json:"server"`
	DB       int    `json:"db"`
	Key      any    `json:"key"`
	Field    string `json:"field,omitempty"`
	NewField string `json:"newField,omitempty"`
	Value    any    `json:"value"`
	Format   string `json:"format,omitempty"`
	Decode   string `json:"decode,omitempty"`
}
