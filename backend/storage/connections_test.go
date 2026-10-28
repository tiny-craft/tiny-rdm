package storage

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"tinyrdm/backend/types"
)

func sortedConnection(name string) types.Connection {
	return types.Connection{ConnectionConfig: types.ConnectionConfig{Name: name}}
}

func sortedGroup(name string, children ...types.Connection) types.Connection {
	return types.Connection{ConnectionConfig: types.ConnectionConfig{Name: name}, Type: "group", Connections: children}
}

func newConnectionSortStorage(t *testing.T) *ConnectionsStorage {
	t.Helper()
	c := &ConnectionsStorage{storage: &localStorage{ConfPath: filepath.Join(t.TempDir(), "connections.yaml")}}
	connections := types.Connections{
		sortedGroup("A", types.Connection{ConnectionConfig: types.ConnectionConfig{
			Name: "cache-10", Addr: "localhost", Port: 6379, Password: "test-secret", LastDB: 3,
			SSH: types.ConnectionSSH{Enable: true, Addr: "jump-host"},
		}}),
		sortedGroup("B"),
		sortedConnection("cache-2"),
	}
	if err := c.saveConnections(connections); err != nil {
		t.Fatal(err)
	}
	return c
}

func TestSaveSortedConnectionPreservesProfilesAndGroupMembership(t *testing.T) {
	c := newConnectionSortStorage(t)
	profile := c.GetConnection("cache-10").ConnectionConfig
	profile.Group = "B"
	if err := c.SaveSortedConnection(types.Connections{
		sortedGroup("B", sortedConnection("cache-2"), sortedConnection("cache-10")),
		sortedGroup("A"),
	}); err != nil {
		t.Fatal(err)
	}
	got := c.GetConnections()
	if got[0].Name != "B" || got[1].Name != "A" || len(got[0].Connections) != 2 || len(got[1].Connections) != 0 {
		t.Fatal("saved order or group membership does not match the requested tree")
	}
	if !reflect.DeepEqual(c.GetConnection("cache-10").ConnectionConfig, profile) {
		t.Fatal("sorting changed the saved connection profile")
	}
}

func TestSaveSortedConnectionRejectsInvalidTreeWithoutChangingFile(t *testing.T) {
	tests := map[string]types.Connections{
		"nested group":         {sortedGroup("B", sortedGroup("A", sortedConnection("cache-10"))), sortedConnection("cache-2")},
		"server with children": {{ConnectionConfig: types.ConnectionConfig{Name: "cache-2"}, Connections: []types.Connection{sortedConnection("cache-10")}}, sortedGroup("A"), sortedGroup("B")},
		"missing connection":   {sortedGroup("A", sortedConnection("cache-10")), sortedGroup("B")},
		"missing empty group":  {sortedGroup("A", sortedConnection("cache-10")), sortedConnection("cache-2")},
		"duplicate connection": {sortedGroup("A", sortedConnection("cache-10")), sortedGroup("B", sortedConnection("cache-10")), sortedConnection("cache-2")},
		"duplicate group":      {sortedGroup("A", sortedConnection("cache-10")), sortedGroup("A"), sortedGroup("B"), sortedConnection("cache-2")},
		"unknown connection":   {sortedGroup("A", sortedConnection("cache-10")), sortedGroup("B"), sortedConnection("missing")},
		"unknown group":        {sortedGroup("A", sortedConnection("cache-10")), sortedGroup("missing"), sortedConnection("cache-2")},
		"empty request":        nil,
	}
	for name, sorted := range tests {
		t.Run(name, func(t *testing.T) {
			c := newConnectionSortStorage(t)
			before, err := os.ReadFile(c.storage.ConfPath)
			if err != nil {
				t.Fatal(err)
			}
			if err := c.SaveSortedConnection(sorted); err == nil {
				t.Fatal("invalid reorder request was accepted")
			}
			after, err := os.ReadFile(c.storage.ConfPath)
			if err != nil {
				t.Fatal(err)
			}
			if string(before) != string(after) {
				t.Fatal("invalid reorder request changed the saved connections")
			}
		})
	}
}

func TestSaveSortedConnectionRejectsStaleListAfterCreatingConnection(t *testing.T) {
	c := newConnectionSortStorage(t)
	stale := c.GetConnections()
	if err := c.CreateConnection(types.ConnectionConfig{Name: "new-connection", Addr: "localhost", Port: 6379}); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(c.storage.ConfPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.SaveSortedConnection(stale); err == nil {
		t.Fatal("stale reorder request was accepted after a connection was created")
	}
	after, err := os.ReadFile(c.storage.ConfPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("stale reorder request changed the saved connections")
	}
}

func TestSaveSortedConnectionEmptyList(t *testing.T) {
	c := &ConnectionsStorage{storage: &localStorage{ConfPath: filepath.Join(t.TempDir(), "connections.yaml")}}
	if err := c.saveConnections(types.Connections{}); err != nil {
		t.Fatal(err)
	}
	if err := c.SaveSortedConnection(types.Connections{}); err != nil {
		t.Fatal(err)
	}
	if len(c.GetConnections()) != 0 {
		t.Fatal("sorting an empty list created connections")
	}
}
