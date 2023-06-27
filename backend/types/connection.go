package types

type ConnectionCategory int

type Connection struct {
	Group         string `json:"group" yaml:"-"`
	Name          string `json:"name" yaml:"name"`
	Addr          string `json:"addr" yaml:"addr"`
	Port          int    `json:"port" yaml:"port"`
	Username      string `json:"username" yaml:"username"`
	Password      string `json:"password" yaml:"password"`
	DefaultFilter string `json:"defaultFilter" yaml:"default_filter"`
	KeySeparator  string `json:"keySeparator" yaml:"key_separator"`
	ConnTimeout   int    `json:"connTimeout" yaml:"conn_timeout"`
	ExecTimeout   int    `json:"execTimeout" yaml:"exec_timeout"`
	MarkColor     string `json:"markColor" yaml:"mark_color"`
}

type ConnectionGroup struct {
	GroupName   string       `json:"groupName" yaml:"group_name"`
	Connections []Connection `json:"connections" yaml:"connections"`
}

type ConnectionDB struct {
	Name    string `json:"name"`
	Keys    int    `json:"keys"`
	Expires int    `json:"expires,omitempty"`
	AvgTTL  int    `json:"avgTtl,omitempty"`
}
