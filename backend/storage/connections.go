package storage

import (
	"errors"
	"gopkg.in/yaml.v3"
	"sync"
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

func (c *ConnectionsStorage) defaultConnections() []types.ConnectionGroup {
	return []types.ConnectionGroup{
		{
			GroupName:   "",
			Connections: []types.Connection{},
		},
	}
}

func (c *ConnectionsStorage) defaultConnectionItem() types.Connection {
	return types.Connection{
		Name:          "",
		Addr:          "127.0.0.1",
		Port:          6379,
		Username:      "",
		Password:      "",
		DefaultFilter: "*",
		KeySeparator:  ":",
		ConnTimeout:   60,
		ExecTimeout:   60,
		MarkColor:     "",
	}
}

func (c *ConnectionsStorage) getConnections() (ret []types.ConnectionGroup) {
	b, err := c.storage.Load()
	if err != nil {
		ret = c.defaultConnections()
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
func (c *ConnectionsStorage) GetConnections() (ret []types.ConnectionGroup) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.getConnections()
}

// GetConnectionsFlat get all store connections from local flat(exclude group level)
func (c *ConnectionsStorage) GetConnectionsFlat() (ret []types.Connection) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	conns := c.getConnections()
	for _, group := range conns {
		for _, conn := range group.Connections {
			ret = append(ret, conn)
		}
	}
	return
}

func (c *ConnectionsStorage) saveConnections(conns []types.ConnectionGroup) error {
	b, err := yaml.Marshal(&conns)
	if err != nil {
		return err
	}
	if err = c.storage.Store(b); err != nil {
		return err
	}
	return nil
}

// UpsertConnection update or insert a connection
func (c *ConnectionsStorage) UpsertConnection(param types.Connection, replace bool) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	conns := c.getConnections()
	groupIndex := -1
	connIndex := -1
	for i, group := range conns {
		for j, conn := range group.Connections {
			// check conflict connection name
			if conn.Name == param.Name {
				if !replace {
					return errors.New("duplicated connection name")
				} else {
					// different group name, should move group
					// remove from current group first
					if group.GroupName != param.Group {
						group.Connections = append(group.Connections[:j], group.Connections[j+1:]...)

						// find new group index
						groupIndex, _ = sliceutil.Find(conns, func(i int) bool {
							return conns[i].GroupName == param.Group
						})
					} else {
						groupIndex = i
						connIndex = j
					}
					break
				}
			}
		}
	}

	if groupIndex >= 0 {
		// group exists
		if connIndex >= 0 {
			// connection exists
			conns[groupIndex].Connections[connIndex] = param
		} else {
			// new connection
			conns[groupIndex].Connections = append(conns[groupIndex].Connections, param)
		}
	} else {
		// new group
		group := types.ConnectionGroup{
			GroupName:   param.Group,
			Connections: []types.Connection{param},
		}
		conns = append(conns, group)
	}

	return c.saveConnections(conns)
}

// RemoveConnection remove special connection
func (c *ConnectionsStorage) RemoveConnection(group, name string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	conns := c.getConnections()
	for i, connGroup := range conns {
		if connGroup.GroupName == group {
			for j, conn := range connGroup.Connections {
				if conn.Name == name {
					connList := conns[i].Connections
					connList = append(connList[:j], connList[j+1:]...)
					conns[i].Connections = connList
					return c.saveConnections(conns)
				}
			}
		}
	}

	return errors.New("no match connection")
}

// UpsertGroup update or insert a group
// When want to create group only, set group == param.name
func (c *ConnectionsStorage) UpsertGroup(group string, param types.ConnectionGroup) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	conns := c.getConnections()
	for i, connGroup := range conns {
		if connGroup.GroupName == group {
			conns[i].GroupName = param.GroupName
			return c.saveConnections(conns)
		}
	}

	// No match group, create one
	connGroup := types.ConnectionGroup{
		GroupName:   param.GroupName,
		Connections: []types.Connection{},
	}
	conns = append(conns, connGroup)
	return c.saveConnections(conns)
}

// RemoveGroup remove special group, include all connections under it
func (c *ConnectionsStorage) RemoveGroup(group string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	conns := c.getConnections()
	for i, connGroup := range conns {
		if connGroup.GroupName == group {
			conns = append(conns[:i], conns[i+1:]...)
			return c.saveConnections(conns)
		}
	}

	return errors.New("no match group")
}
