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
	Format       string `json:"format,omitempty"`
	Decode       string `json:"decode,omitempty"`
	MatchPattern string `json:"matchPattern,omitempty"`
	Reset        bool   `json:"reset"`
	Full         bool   `json:"full"`
}

type KeyDetail struct {
	Value  any    `json:"value"`
	Length int64  `json:"length,omitempty"`
	Format string `json:"format,omitempty"`
	Decode string `json:"decode,omitempty"`
	End    bool   `json:"end"`
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

type SetListParam struct {
	Server    string `json:"server"`
	DB        int    `json:"db"`
	Key       any    `json:"key"`
	Index     int64  `json:"index"`
	Value     any    `json:"value"`
	Format    string `json:"format,omitempty"`
	Decode    string `json:"decode,omitempty"`
	RetFormat string `json:"retFormat,omitempty"`
	RetDecode string `json:"retDecode,omitempty"`
}

type SetHashParam struct {
	Server    string `json:"server"`
	DB        int    `json:"db"`
	Key       any    `json:"key"`
	Field     string `json:"field,omitempty"`
	NewField  string `json:"newField,omitempty"`
	Value     any    `json:"value"`
	Format    string `json:"format,omitempty"`
	Decode    string `json:"decode,omitempty"`
	RetFormat string `json:"retFormat,omitempty"`
	RetDecode string `json:"retDecode,omitempty"`
}

type SetSetParam struct {
	Server    string `json:"server"`
	DB        int    `json:"db"`
	Key       any    `json:"key"`
	Value     any    `json:"value"`
	NewValue  any    `json:"newValue"`
	Format    string `json:"format,omitempty"`
	Decode    string `json:"decode,omitempty"`
	RetFormat string `json:"retFormat,omitempty"`
	RetDecode string `json:"retDecode,omitempty"`
}

type SetZSetParam struct {
	Server    string  `json:"server"`
	DB        int     `json:"db"`
	Key       any     `json:"key"`
	Value     any     `json:"value"`
	NewValue  any     `json:"newValue"`
	Score     float64 `json:"score"`
	Format    string  `json:"format,omitempty"`
	Decode    string  `json:"decode,omitempty"`
	RetFormat string  `json:"retFormat,omitempty"`
	RetDecode string  `json:"retDecode,omitempty"`
}
