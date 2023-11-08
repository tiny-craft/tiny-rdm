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
