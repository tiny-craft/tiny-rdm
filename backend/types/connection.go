package types

type ConnectionCategory int

type ConnectionConfig struct {
	Name          string `json:"name" yaml:"name"`
	Group         string `json:"group,omitempty" yaml:"-"`
	Addr          string `json:"addr,omitempty" yaml:"addr,omitempty"`
	Port          int    `json:"port,omitempty" yaml:"port,omitempty"`
	Username      string `json:"username,omitempty" yaml:"username,omitempty"`
	Password      string `json:"password,omitempty" yaml:"password,omitempty"`
	DefaultFilter string `json:"defaultFilter,omitempty" yaml:"default_filter,omitempty"`
	KeySeparator  string `json:"keySeparator,omitempty" yaml:"key_separator,omitempty"`
	ConnTimeout   int    `json:"connTimeout,omitempty" yaml:"conn_timeout,omitempty"`
	ExecTimeout   int    `json:"execTimeout,omitempty" yaml:"exec_timeout,omitempty"`
	MarkColor     string `json:"markColor,omitempty" yaml:"mark_color,omitempty"`
	SafeLink      int    `json:"safeLink,omitempty" yaml:"safe_link,omitempty"`
	SshAddr       string `json:"sshAddr,omitempty" yaml:"ssh_addr,omitempty"`
	SshPort       int    `json:"sshPort,omitempty" yaml:"ssh_port,omitempty"`
	SshUser       string `json:"sshUser,omitempty" yaml:"ssh_user,omitempty"`
	SshAuth       int    `json:"sshAuth,omitempty" yaml:"ssh_auth,omitempty"`
	SshKeyPath    string `json:"sshKeyPath,omitempty" yaml:"ssh_key_path,omitempty"`
	SshKeyPwd     string `json:"sshKeyPwd,omitempty" yaml:"ssh_key_pwd,omitempty"`
	SshPassword   string `json:"sshPassword,omitempty" yaml:"ssh_password,omitempty"`
}

type Connection struct {
	ConnectionConfig `json:",inline" yaml:",inline"`
	Type             string       `json:"type,omitempty" yaml:"type,omitempty"`
	Connections      []Connection `json:"connections,omitempty" yaml:"connections,omitempty"`
}

type Connections []Connection

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
