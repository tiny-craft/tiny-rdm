package storage

import (
	"errors"
	"gopkg.in/yaml.v3"
	"sync"
	"tinyrdm/backend/consts"
	"tinyrdm/backend/types"
	sliceutil "tinyrdm/backend/utils/slice"
)

type ConnectionsStorage struct {
	storage *localStorage
	mutex   sync.Mutex
}

func NewConnections() *ConnectionsStorage {
	return &ConnectionsStorage{
		storage: NewLocalStore("connections.yaml"),
	}
}

func (c *ConnectionsStorage) defaultConnections() types.Connections {
	return types.Connections{}
}

func (c *ConnectionsStorage) defaultConnectionItem() types.ConnectionConfig {
	return types.ConnectionConfig{
		Name:          "",
		Addr:          "127.0.0.1",
		Port:          6379,
		Username:      "",
		Password:      "",
		DefaultFilter: "*",
		KeySeparator:  ":",
		ConnTimeout:   60,
		ExecTimeout:   60,
		DBFilterType:  "none",
		DBFilterList:  []int{},
		LoadSize:      consts.DEFAULT_LOAD_SIZE,
		MarkColor:     "",
		Sentinel: types.ConnectionSentinel{
			Master: "mymaster",
		},
	}
}

func (c *ConnectionsStorage) getConnections() (ret types.Connections) {
	b, err := c.storage.Load()
	ret = c.defaultConnections()
	if err != nil {
		return
	}

	if err = yaml.Unmarshal(b, &ret); err != nil {
		ret = c.defaultConnections()
		return
	}
	if len(ret) <= 0 {
		ret = c.defaultConnections()
	}
	//if !sliceutil.AnyMatch(ret, func(i int) bool {
	//	return ret[i].GroupName == ""
	//}) {
	//	ret = append(ret, c.defaultConnections()...)
	//}
	return
}

// GetConnections get all store connections from local
func (c *ConnectionsStorage) GetConnections() (ret types.Connections) {
	return c.getConnections()
}

// GetConnectionsFlat get all store connections from local flat(exclude group level)
func (c *ConnectionsStorage) GetConnectionsFlat() (ret types.Connections) {
	conns := c.getConnections()
	for _, conn := range conns {
		if conn.Type == "group" {
			ret = append(ret, conn.Connections...)
		} else {
			ret = append(ret, conn)
		}
	}
	return
}

// GetConnection get connection by name
func (c *ConnectionsStorage) GetConnection(name string) *types.Connection {
	conns := c.getConnections()

	var findConn func(string, string, types.Connections) *types.Connection
	findConn = func(name, groupName string, conns types.Connections) *types.Connection {
		for i, conn := range conns {
			if conn.Type != "group" {
				if conn.Name == name {
					conns[i].Group = groupName
					return &conns[i]
				}
			} else {
				if ret := findConn(name, conn.Name, conn.Connections); ret != nil {
					return ret
				}
			}
		}
		return nil
	}

	return findConn(name, "", conns)
}

// GetGroup get one connection group by name
func (c *ConnectionsStorage) GetGroup(name string) *types.Connection {
	conns := c.getConnections()

	for i, conn := range conns {
		if conn.Type == "group" && conn.Name == name {
			return &conns[i]
		}
	}
	return nil
}

func (c *ConnectionsStorage) saveConnections(conns types.Connections) error {
	b, err := yaml.Marshal(&conns)
	if err != nil {
		return err
	}
	if err = c.storage.Store(b); err != nil {
		return err
	}
	return nil
}

// CreateConnection create new connection
func (c *ConnectionsStorage) CreateConnection(param types.ConnectionConfig) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	conn := c.GetConnection(param.Name)
	if conn != nil {
		return errors.New("duplicated connection name")
	}

	conns := c.getConnections()
	var group *types.Connection
	if len(param.Group) > 0 {
		for i, conn := range conns {
			if conn.Type == "group" && conn.Name == param.Group {
				group = &conns[i]
				break
			}
		}
	}
	if group != nil {
		group.Connections = append(group.Connections, types.Connection{
			ConnectionConfig: param,
		})
	} else {
		if len(param.Group) > 0 {
			// no group matched, create new group
			conns = append(conns, types.Connection{
				Type: "group",
				Connections: types.Connections{
					types.Connection{
						ConnectionConfig: param,
					},
				},
			})
		} else {
			conns = append(conns, types.Connection{
				ConnectionConfig: param,
			})
		}
	}

	return c.saveConnections(conns)
}

// UpdateConnection update existing connection by name
func (c *ConnectionsStorage) UpdateConnection(name string, param types.ConnectionConfig) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	conns := c.getConnections()
	var updated bool
	var retrieve func(types.Connections, string, types.ConnectionConfig) error
	retrieve = func(conns types.Connections, name string, param types.ConnectionConfig) error {
		for i, conn := range conns {
			if conn.Type != "group" {
				if name != param.Name && conn.Name == param.Name {
					return errors.New("duplicated connection name")
				} else if conn.Name == name && !updated {
					conns[i] = types.Connection{
						ConnectionConfig: param,
					}
					updated = true
				}
			} else {
				if err := retrieve(conn.Connections, name, param); err != nil {
					return err
				}
			}
		}
		return nil
	}

	err := retrieve(conns, name, param)
	if err != nil {
		return err
	}
	if !updated {
		return errors.New("connection not found")
	}

	return c.saveConnections(conns)
}

// DeleteConnection remove special connection
func (c *ConnectionsStorage) DeleteConnection(name string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	conns := c.getConnections()
	var updated bool
	for i, conn := range conns {
		if conn.Type == "group" {
			for j, subConn := range conn.Connections {
				if subConn.Name == name {
					conns[i].Connections = append(conns[i].Connections[:j], conns[i].Connections[j+1:]...)
					updated = true
					break
				}
			}
		} else if conn.Name == name {
			conns = append(conns[:i], conns[i+1:]...)
			updated = true
			break
		}
		if updated {
			break
		}
	}
	if !updated {
		return errors.New("no match connection")
	}
	return c.saveConnections(conns)
}

// SaveSortedConnection save connection after sort
func (c *ConnectionsStorage) SaveSortedConnection(sortedConns types.Connections) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	conns := c.GetConnectionsFlat()
	takeConn := func(name string) (types.Connection, bool) {
		idx, ok := sliceutil.Find(conns, func(i int) bool {
			return conns[i].Name == name
		})
		if ok {
			ret := conns[idx]
			conns = append(conns[:idx], conns[idx+1:]...)
			return ret, true
		}
		return types.Connection{}, false
	}
	var replaceConn func(connections types.Connections) types.Connections
	replaceConn = func(cons types.Connections) types.Connections {
		var newConns types.Connections
		for _, conn := range cons {
			if conn.Type == "group" {
				newConns = append(newConns, types.Connection{
					ConnectionConfig: types.ConnectionConfig{
						Name: conn.Name,
					},
					Type:        "group",
					Connections: replaceConn(conn.Connections),
				})
			} else {
				if foundConn, ok := takeConn(conn.Name); ok {
					newConns = append(newConns, foundConn)
				}
			}
		}
		return newConns
	}
	conns = replaceConn(sortedConns)
	return c.saveConnections(conns)
}

// CreateGroup create a new group
func (c *ConnectionsStorage) CreateGroup(name string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	conns := c.getConnections()
	for _, conn := range conns {
		if conn.Type == "group" && conn.Name == name {
			return errors.New("duplicated group name")
		}
	}

	conns = append(conns, types.Connection{
		ConnectionConfig: types.ConnectionConfig{
			Name: name,
		},
		Type: "group",
	})
	return c.saveConnections(conns)
}

// RenameGroup rename group
func (c *ConnectionsStorage) RenameGroup(name, newName string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	groupIndex := -1
	conns := c.getConnections()
	for i, conn := range conns {
		if conn.Type == "group" {
			if conn.Name == newName {
				return errors.New("duplicated group name")
			} else if conn.Name == name {
				groupIndex = i
			}
		}
	}

	if groupIndex == -1 {
		return errors.New("group not found")
	}

	conns[groupIndex].Name = newName
	return c.saveConnections(conns)
}

// DeleteGroup remove specified group, include all connections under it
func (c *ConnectionsStorage) DeleteGroup(group string, includeConnection bool) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	conns := c.getConnections()
	for i, conn := range conns {
		if conn.Type == "group" && conn.Name == group {
			conns = append(conns[:i], conns[i+1:]...)
			if includeConnection {
				conns = append(conns, conn.Connections...)
			}
			return c.saveConnections(conns)
		}
	}
	return errors.New("group not found")
}
