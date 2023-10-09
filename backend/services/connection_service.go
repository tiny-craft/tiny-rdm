package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	. "tinyrdm/backend/storage"
	"tinyrdm/backend/types"
	"tinyrdm/backend/utils/coll"
	maputil "tinyrdm/backend/utils/map"
	redis2 "tinyrdm/backend/utils/redis"
	sliceutil "tinyrdm/backend/utils/slice"
	strutil "tinyrdm/backend/utils/string"
)

type cmdHistoryItem struct {
	Timestamp int64  `json:"timestamp"`
	Server    string `json:"server"`
	Cmd       string `json:"cmd"`
	Cost      int64  `json:"cost"`
}

type connectionService struct {
	ctx        context.Context
	conns      *ConnectionsStorage
	connMap    map[string]connectionItem
	cmdHistory []cmdHistoryItem
}

type connectionItem struct {
	rdb        *redis.Client
	ctx        context.Context
	cancelFunc context.CancelFunc
}

type keyItem struct {
	Type string `json:"t"`
}

var connection *connectionService
var onceConnection sync.Once

func Connection() *connectionService {
	if connection == nil {
		onceConnection.Do(func() {
			connection = &connectionService{
				conns:   NewConnections(),
				connMap: map[string]connectionItem{},
			}
		})
	}
	return connection
}

func (c *connectionService) Start(ctx context.Context) {
	c.ctx = ctx
}

func (c *connectionService) Stop(ctx context.Context) {
	for _, item := range c.connMap {
		if item.rdb != nil {
			item.cancelFunc()
			item.rdb.Close()
		}
	}
	c.connMap = map[string]connectionItem{}
}

func (c *connectionService) buildOption(config types.ConnectionConfig) (*redis.Options, error) {
	var sshClient *ssh.Client
	if config.SSH.Enable {
		sshConfig := &ssh.ClientConfig{
			User:            config.SSH.Username,
			Auth:            []ssh.AuthMethod{ssh.Password(config.SSH.Password)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         time.Duration(config.ConnTimeout) * time.Second,
		}
		switch config.SSH.LoginType {
		case "pwd":
			sshConfig.Auth = []ssh.AuthMethod{ssh.Password(config.SSH.Password)}
		case "pkfile":
			key, err := os.ReadFile(config.SSH.PKFile)
			if err != nil {
				return nil, err
			}
			var signer ssh.Signer
			if len(config.SSH.Passphrase) > 0 {
				signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(config.SSH.Passphrase))
			} else {
				signer, err = ssh.ParsePrivateKey(key)
			}
			if err != nil {
				return nil, err
			}
			sshConfig.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
		default:
			return nil, errors.New("invalid login type")
		}

		var err error
		sshClient, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", config.SSH.Addr, config.SSH.Port), sshConfig)
		if err != nil {
			return nil, err
		}
	}

	option := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Addr, config.Port),
		Username:     config.Username,
		Password:     config.Password,
		DialTimeout:  time.Duration(config.ConnTimeout) * time.Second,
		ReadTimeout:  time.Duration(config.ExecTimeout) * time.Second,
		WriteTimeout: time.Duration(config.ExecTimeout) * time.Second,
	}
	if sshClient != nil {
		option.Dialer = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return sshClient.Dial(network, addr)
		}
		option.ReadTimeout = -2
		option.WriteTimeout = -2
	}
	return option, nil
}

func (c *connectionService) createRedisClient(config types.ConnectionConfig) (*redis.Client, error) {
	option, err := c.buildOption(config)
	if err != nil {
		return nil, err
	}

	if config.Sentinel.Enable {
		sentinel := redis.NewSentinelClient(option)
		defer sentinel.Close()

		var addr []string
		addr, err = sentinel.GetMasterAddrByName(c.ctx, config.Sentinel.Master).Result()
		if err != nil {
			return nil, err
		}
		if len(addr) < 2 {
			return nil, errors.New("cannot get master address")
		}
		option.Addr = fmt.Sprintf("%s:%s", addr[0], addr[1])
		option.Username = config.Sentinel.Username
		option.Password = config.Sentinel.Password
	}

	rdb := redis.NewClient(option)
	return rdb, nil
}

// ListSentinelMasters list all master info by sentinel
func (c *connectionService) ListSentinelMasters(config types.ConnectionConfig) (resp types.JSResp) {
	option, err := c.buildOption(config)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	if option.DialTimeout > 0 {
		option.DialTimeout = 10 * time.Second
	}
	sentinel := redis.NewSentinelClient(option)
	defer sentinel.Close()

	var retInfo []map[string]string
	masterInfos, err := sentinel.Masters(c.ctx).Result()
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	for _, info := range masterInfos {
		if infoMap, ok := info.(map[any]any); ok {
			retInfo = append(retInfo, map[string]string{
				"name": infoMap["name"].(string),
				"addr": fmt.Sprintf("%s:%s", infoMap["ip"].(string), infoMap["port"].(string)),
			})
		}
	}

	resp.Data = retInfo
	resp.Success = true
	return
}

func (c *connectionService) TestConnection(config types.ConnectionConfig) (resp types.JSResp) {
	rdb, err := c.createRedisClient(config)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	defer rdb.Close()

	if _, err = rdb.Ping(c.ctx).Result(); err != nil && err != redis.Nil {
		resp.Msg = err.Error()
	} else {
		resp.Success = true
	}
	return
}

// ListConnection list all saved connection in local profile
func (c *connectionService) ListConnection() (resp types.JSResp) {
	resp.Success = true
	resp.Data = c.conns.GetConnections()
	return
}

// GetConnection get connection profile by name
func (c *connectionService) GetConnection(name string) (resp types.JSResp) {
	conn := c.conns.GetConnection(name)
	resp.Success = conn != nil
	resp.Data = conn
	return
}

// SaveConnection save connection config to local profile
func (c *connectionService) SaveConnection(name string, param types.ConnectionConfig) (resp types.JSResp) {
	var err error
	if strings.ContainsAny(param.Name, "/") {
		err = errors.New("connection name contains illegal characters")
	} else {
		if len(name) > 0 {
			// update connection
			err = c.conns.UpdateConnection(name, param)
		} else {
			err = c.conns.CreateConnection(param)
		}
	}
	if err != nil {
		resp.Msg = err.Error()
	} else {
		resp.Success = true
	}
	return
}

// DeleteConnection remove connection by name
func (c *connectionService) DeleteConnection(name string) (resp types.JSResp) {
	err := c.conns.DeleteConnection(name)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Success = true
	return
}

// SaveSortedConnection save sorted connection after drag
func (c *connectionService) SaveSortedConnection(sortedConns types.Connections) (resp types.JSResp) {
	err := c.conns.SaveSortedConnection(sortedConns)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Success = true
	return
}

// SelectKeyFile open file dialog to select a private key file
func (c *connectionService) SelectKeyFile(title string) (resp types.JSResp) {
	filepath, err := runtime.OpenFileDialog(c.ctx, runtime.OpenDialogOptions{
		Title:           title,
		ShowHiddenFiles: true,
	})
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Success = true
	resp.Data = map[string]any{
		"path": filepath,
	}
	return
}

// CreateGroup create a new group
func (c *connectionService) CreateGroup(name string) (resp types.JSResp) {
	err := c.conns.CreateGroup(name)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Success = true
	return
}

// RenameGroup rename group
func (c *connectionService) RenameGroup(name, newName string) (resp types.JSResp) {
	err := c.conns.RenameGroup(name, newName)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Success = true
	return
}

// DeleteGroup remove a group by name
func (c *connectionService) DeleteGroup(name string, includeConn bool) (resp types.JSResp) {
	err := c.conns.DeleteGroup(name, includeConn)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Success = true
	return
}

// OpenConnection open redis server connection
func (c *connectionService) OpenConnection(name string) (resp types.JSResp) {
	rdb, ctx, err := c.getRedisClient(name, 0)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	// get connection config
	selConn := c.conns.GetConnection(name)

	var totaldb int
	if selConn.DBFilterType == "" || selConn.DBFilterType == "none" {
		// get total databases
		if config, err := rdb.ConfigGet(ctx, "databases").Result(); err == nil {
			if total, err := strconv.Atoi(config["databases"]); err == nil {
				totaldb = total
			}
		}
	}

	// get database info
	res, err := rdb.Info(ctx, "keyspace").Result()
	if err != nil {
		resp.Msg = "get server info fail:" + err.Error()
		return
	}
	// parse all db, response content like below
	var dbs []types.ConnectionDB
	info := c.parseInfo(res)
	if totaldb <= 0 {
		// cannot retrieve the database count by "CONFIG GET databases", try to get max index from keyspace
		keyspace := info["Keyspace"]
		var db, maxDB int
		for dbName := range keyspace {
			if db, err = strconv.Atoi(strings.TrimLeft(dbName, "db")); err == nil {
				if maxDB < db {
					maxDB = db
				}
			}
		}
		totaldb = maxDB + 1
	}

	queryDB := func(idx int) types.ConnectionDB {
		dbName := "db" + strconv.Itoa(idx)
		dbInfoStr := info["Keyspace"][dbName]
		if len(dbInfoStr) > 0 {
			dbInfo := c.parseDBItemInfo(dbInfoStr)
			return types.ConnectionDB{
				Name:    dbName,
				Index:   idx,
				Keys:    dbInfo["keys"],
				Expires: dbInfo["expires"],
				AvgTTL:  dbInfo["avg_ttl"],
			}
		} else {
			return types.ConnectionDB{
				Name:  dbName,
				Index: idx,
			}
		}
	}
	switch selConn.DBFilterType {
	case "show":
		filterList := sliceutil.Unique(selConn.DBFilterList)
		for _, idx := range filterList {
			dbs = append(dbs, queryDB(idx))
		}
	case "hide":
		hiddenList := coll.NewSet(selConn.DBFilterList...)
		for idx := 0; idx < totaldb; idx++ {
			if !hiddenList.Contains(idx) {
				dbs = append(dbs, queryDB(idx))
			}
		}
	default:
		for idx := 0; idx < totaldb; idx++ {
			dbs = append(dbs, queryDB(idx))
		}
	}

	resp.Success = true
	resp.Data = map[string]any{
		"db": dbs,
	}
	return
}

// CloseConnection close redis server connection
func (c *connectionService) CloseConnection(name string) (resp types.JSResp) {
	item, ok := c.connMap[name]
	if ok {
		delete(c.connMap, name)
		if item.rdb != nil {
			item.cancelFunc()
			item.rdb.Close()
		}
	}
	resp.Success = true
	return
}

// get redis client from local cache or create a new open
// if db >= 0, will also switch to db index
func (c *connectionService) getRedisClient(connName string, db int) (*redis.Client, context.Context, error) {
	item, ok := c.connMap[connName]
	var rdb *redis.Client
	var ctx context.Context
	if ok {
		rdb, ctx = item.rdb, item.ctx
	} else {
		selConn := c.conns.GetConnection(connName)
		if selConn == nil {
			return nil, nil, fmt.Errorf("no match connection \"%s\"", connName)
		}

		var err error
		rdb, err = c.createRedisClient(selConn.ConnectionConfig)
		if err != nil {
			return nil, nil, fmt.Errorf("create conenction error: %s", err.Error())
		}
		rdb.AddHook(redis2.NewHook(connName, func(cmd string, cost int64) {
			now := time.Now()
			//last := strings.LastIndex(cmd, ":")
			//if last != -1 {
			//	cmd = cmd[:last]
			//}
			c.cmdHistory = append(c.cmdHistory, cmdHistoryItem{
				Timestamp: now.UnixMilli(),
				Server:    connName,
				Cmd:       cmd,
				Cost:      cost,
			})
		}))

		if _, err = rdb.Ping(c.ctx).Result(); err != nil && err != redis.Nil {
			return nil, nil, errors.New("can not connect to redis server:" + err.Error())
		}
		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithCancel(c.ctx)
		c.connMap[connName] = connectionItem{
			rdb:        rdb,
			ctx:        ctx,
			cancelFunc: cancelFunc,
		}
	}

	if db >= 0 {
		if err := rdb.Do(ctx, "select", strconv.Itoa(db)).Err(); err != nil {
			return nil, nil, err
		}
	}
	return rdb, ctx, nil
}

// parse command response content which use "redis info"
// # Keyspace\r\ndb0:keys=2,expires=1,avg_ttl=1877111749\r\ndb1:keys=33,expires=0,avg_ttl=0\r\ndb3:keys=17,expires=0,avg_ttl=0\r\ndb5:keys=3,expires=0,avg_ttl=0\r\n
func (c *connectionService) parseInfo(info string) map[string]map[string]string {
	parsedInfo := map[string]map[string]string{}
	lines := strings.Split(info, "\r\n")
	if len(lines) > 0 {
		var subInfo map[string]string
		for _, line := range lines {
			if strings.HasPrefix(line, "#") {
				subInfo = map[string]string{}
				parsedInfo[strings.TrimSpace(strings.TrimLeft(line, "#"))] = subInfo
			} else {
				items := strings.SplitN(line, ":", 2)
				if len(items) < 2 {
					continue
				}
				subInfo[items[0]] = items[1]
			}
		}
	}
	return parsedInfo
}

// parse db item value, content format like below
// keys=2,expires=1,avg_ttl=1877111749
func (c *connectionService) parseDBItemInfo(info string) map[string]int {
	ret := map[string]int{}
	items := strings.Split(info, ",")
	for _, item := range items {
		kv := strings.SplitN(item, "=", 2)
		if len(kv) > 1 {
			ret[kv[0]], _ = strconv.Atoi(kv[1])
		}
	}
	return ret
}

// ServerInfo get server info
func (c *connectionService) ServerInfo(name string) (resp types.JSResp) {
	rdb, ctx, err := c.getRedisClient(name, 0)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	// get database info
	res, err := rdb.Info(ctx).Result()
	if err != nil {
		resp.Msg = "get server info fail:" + err.Error()
		return
	}

	resp.Success = true
	resp.Data = c.parseInfo(res)
	return
}

// OpenDatabase open select database, and list all keys
// @param path contain connection name and db name
func (c *connectionService) OpenDatabase(connName string, db int, match string, keyType string) (resp types.JSResp) {
	return c.ScanKeys(connName, db, match, keyType)
}

// ScanKeys scan all keys
func (c *connectionService) ScanKeys(connName string, db int, match, keyType string) (resp types.JSResp) {
	rdb, ctx, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	filterType := len(keyType) > 0

	var keys []string
	//keys := map[string]keyItem{}
	var cursor uint64
	for {
		var loadedKey []string
		if filterType {
			loadedKey, cursor, err = rdb.ScanType(ctx, cursor, match, 10000, keyType).Result()
		} else {
			loadedKey, cursor, err = rdb.Scan(ctx, cursor, match, 10000).Result()
		}
		if err != nil {
			resp.Msg = err.Error()
			return
		}
		keys = append(keys, loadedKey...)
		//for _, k := range loadedKey {
		//	//t, _ := rdb.Type(ctx, k).Result()
		//	keys[k] = keyItem{Type: "t"}
		//}
		//keys = append(keys, loadedKey...)
		// no more loadedKey
		if cursor == 0 {
			break
		}
	}

	resp.Success = true
	resp.Data = map[string]any{
		"keys": keys,
	}
	return
}

// GetKeyValue get value by key
func (c *connectionService) GetKeyValue(connName string, db int, key, viewAs string) (resp types.JSResp) {
	rdb, ctx, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	var keyType string
	var dur time.Duration
	keyType, err = rdb.Type(ctx, key).Result()
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	if keyType == "none" {
		resp.Msg = "key not exists"
		return
	}

	var ttl int64
	if dur, err = rdb.TTL(ctx, key).Result(); err != nil {
		ttl = -1
	} else {
		if dur < 0 {
			ttl = -1
		} else {
			ttl = int64(dur.Seconds())
		}
	}

	var value any
	var size int64
	var cursor uint64
	switch strings.ToLower(keyType) {
	case "string":
		var str string
		str, err = rdb.Get(ctx, key).Result()
		value, viewAs = strutil.ConvertTo(str, viewAs)
		size, _ = rdb.StrLen(ctx, key).Result()
	case "list":
		value, err = rdb.LRange(ctx, key, 0, -1).Result()
		size, _ = rdb.LLen(ctx, key).Result()
	case "hash":
		//value, err = rdb.HGetAll(ctx, key).Result()
		items := map[string]string{}
		for {
			var loadedVal []string
			loadedVal, cursor, err = rdb.HScan(ctx, key, cursor, "*", 10000).Result()
			if err != nil {
				resp.Msg = err.Error()
				return
			}
			for i := 0; i < len(loadedVal); i += 2 {
				items[loadedVal[i]] = loadedVal[i+1]
			}
			if cursor == 0 {
				break
			}
		}
		value = items
		size, _ = rdb.HLen(ctx, key).Result()
	case "set":
		//value, err = rdb.SMembers(ctx, key).Result()
		items := []string{}
		for {
			var loadedKey []string
			loadedKey, cursor, err = rdb.SScan(ctx, key, cursor, "*", 10000).Result()
			if err != nil {
				resp.Msg = err.Error()
				return
			}
			items = append(items, loadedKey...)
			if cursor == 0 {
				break
			}
		}
		value = items
		size, _ = rdb.SCard(ctx, key).Result()
	case "zset":
		//value, err = rdb.ZRangeWithScores(ctx, key, 0, -1).Result()
		var items []types.ZSetItem
		for {
			var loadedVal []string
			loadedVal, cursor, err = rdb.ZScan(ctx, key, cursor, "*", 10000).Result()
			if err != nil {
				resp.Msg = err.Error()
				return
			}
			var score float64
			for i := 0; i < len(loadedVal); i += 2 {
				if score, err = strconv.ParseFloat(loadedVal[i+1], 64); err == nil {
					items = append(items, types.ZSetItem{
						Value: loadedVal[i],
						Score: score,
					})
				}
			}
			if cursor == 0 {
				break
			}
		}
		value = items
		size, _ = rdb.ZCard(ctx, key).Result()
	case "stream":
		var msgs []redis.XMessage
		items := []types.StreamItem{}
		msgs, err = rdb.XRevRange(ctx, key, "+", "-").Result()
		if err != nil {
			resp.Msg = err.Error()
			return
		}
		for _, msg := range msgs {
			items = append(items, types.StreamItem{
				ID:    msg.ID,
				Value: msg.Values,
			})
		}
		value = items
		size, _ = rdb.XLen(ctx, key).Result()
	}
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Success = true
	resp.Data = map[string]any{
		"type":   keyType,
		"ttl":    ttl,
		"value":  value,
		"size":   size,
		"viewAs": viewAs,
	}
	return
}

// SetKeyValue set value by key
// @param ttl <= 0 means keep current ttl
func (c *connectionService) SetKeyValue(connName string, db int, key, keyType string, value any, ttl int64, viewAs string) (resp types.JSResp) {
	rdb, ctx, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	var expiration time.Duration
	if ttl < 0 {
		if expiration, err = rdb.PTTL(ctx, key).Result(); err != nil {
			expiration = redis.KeepTTL
		}
	} else {
		expiration = time.Duration(ttl) * time.Second
	}
	switch strings.ToLower(keyType) {
	case "string":
		if str, ok := value.(string); !ok {
			resp.Msg = "invalid string value"
			return
		} else {
			var saveStr string
			if saveStr, err = strutil.SaveAs(str, viewAs); err != nil {
				resp.Msg = fmt.Sprintf(`save to "%s" type fail: %s`, viewAs, err.Error())
				return
			}
			_, err = rdb.Set(ctx, key, saveStr, 0).Result()
			// set expiration lonely, not "keepttl"
			if err == nil && expiration > 0 {
				rdb.Expire(ctx, key, expiration)
			}
		}
	case "list":
		if strs, ok := value.([]any); !ok {
			resp.Msg = "invalid list value"
			return
		} else {
			err = rdb.LPush(ctx, key, strs...).Err()
			if err == nil && expiration > 0 {
				rdb.Expire(ctx, key, expiration)
			}
		}
	case "hash":
		if strs, ok := value.([]any); !ok {
			resp.Msg = "invalid hash value"
			return
		} else {
			total := len(strs)
			if total > 1 {
				_, err = rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
					for i := 0; i < total; i += 2 {
						pipe.HSet(ctx, key, strs[i], strs[i+1])
					}
					if expiration > 0 {
						pipe.Expire(ctx, key, expiration)
					}
					return nil
				})
			}
		}
	case "set":
		if strs, ok := value.([]any); !ok || len(strs) <= 0 {
			resp.Msg = "invalid set value"
			return
		} else {
			if len(strs) > 0 {
				err = rdb.SAdd(ctx, key, strs...).Err()
				if err == nil && expiration > 0 {
					rdb.Expire(ctx, key, expiration)
				}
			}
		}
	case "zset":
		if strs, ok := value.([]any); !ok || len(strs) <= 0 {
			resp.Msg = "invalid zset value"
			return
		} else {
			if len(strs) > 1 {
				var members []redis.Z
				for i := 0; i < len(strs); i += 2 {
					score, _ := strconv.ParseFloat(strs[i+1].(string), 64)
					members = append(members, redis.Z{
						Score:  score,
						Member: strs[i],
					})
				}
				err = rdb.ZAdd(ctx, key, members...).Err()
				if err == nil && expiration > 0 {
					rdb.Expire(ctx, key, expiration)
				}
			}
		}
	case "stream":
		if strs, ok := value.([]any); !ok {
			resp.Msg = "invalid stream value"
			return
		} else {
			if len(strs) > 2 {
				err = rdb.XAdd(ctx, &redis.XAddArgs{
					Stream: key,
					ID:     strs[0].(string),
					Values: strs[1:],
				}).Err()
				if err == nil && expiration > 0 {
					rdb.Expire(ctx, key, expiration)
				}
			}
		}
	}

	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Success = true
	resp.Data = map[string]any{
		"value": value,
	}
	return
}

// SetHashValue set hash field
func (c *connectionService) SetHashValue(connName string, db int, key, field, newField, value string) (resp types.JSResp) {
	rdb, ctx, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	var removedField []string
	updatedField := map[string]string{}
	if len(field) <= 0 {
		// old filed is empty, add new field
		_, err = rdb.HSet(ctx, key, newField, value).Result()
		updatedField[newField] = value
	} else if len(newField) <= 0 {
		// new field is empty, delete old field
		_, err = rdb.HDel(ctx, key, field, value).Result()
		removedField = append(removedField, field)
	} else if field == newField {
		// replace field
		_, err = rdb.HSet(ctx, key, newField, value).Result()
		updatedField[newField] = value
	} else {
		// remove old field and add new field
		if _, err = rdb.HDel(ctx, key, field).Result(); err != nil {
			resp.Msg = err.Error()
			return
		}
		_, err = rdb.HSet(ctx, key, newField, value).Result()
		removedField = append(removedField, field)
		updatedField[newField] = value
	}
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	resp.Success = true
	resp.Data = map[string]any{
		"removed": removedField,
		"updated": updatedField,
	}
	return
}

// AddHashField add or update hash field
func (c *connectionService) AddHashField(connName string, db int, key string, action int, fieldItems []any) (resp types.JSResp) {
	rdb, ctx, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	updated := map[string]any{}
	switch action {
	case 1:
		// ignore duplicated fields
		for i := 0; i < len(fieldItems); i += 2 {
			_, err = rdb.HSetNX(ctx, key, fieldItems[i].(string), fieldItems[i+1]).Result()
			if err == nil {
				updated[fieldItems[i].(string)] = fieldItems[i+1]
			}
		}
	default:
		// overwrite duplicated fields
		total := len(fieldItems)
		if total > 1 {
			_, err = rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
				for i := 0; i < total; i += 2 {
					rdb.HSet(ctx, key, fieldItems[i], fieldItems[i+1])
				}
				return nil
			})
			for i := 0; i < total; i += 2 {
				updated[fieldItems[i].(string)] = fieldItems[i+1]
			}
		}
	}
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	resp.Success = true
	resp.Data = map[string]any{
		"updated": updated,
	}
	return
}

// AddListItem add item to list or remove from it
func (c *connectionService) AddListItem(connName string, db int, key string, action int, items []any) (resp types.JSResp) {
	rdb, ctx, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	var leftPush, rightPush []any
	switch action {
	case 0:
		// push to head
		_, err = rdb.LPush(ctx, key, items...).Result()
		leftPush = append(leftPush, items...)
	default:
		// append to tail
		_, err = rdb.RPush(ctx, key, items...).Result()
		rightPush = append(rightPush, items...)
	}
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	resp.Success = true
	resp.Data = map[string]any{
		"left":  leftPush,
		"right": rightPush,
	}
	return
}

// SetListItem update or remove list item by index
func (c *connectionService) SetListItem(connName string, db int, key string, index int64, value string) (resp types.JSResp) {
	rdb, ctx, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	var removed []int64
	updated := map[int64]string{}
	if len(value) <= 0 {
		// remove from list
		err = rdb.LSet(ctx, key, index, "---VALUE_REMOVED_BY_TINY_RDM---").Err()
		if err != nil {
			resp.Msg = err.Error()
			return
		}

		err = rdb.LRem(ctx, key, 1, "---VALUE_REMOVED_BY_TINY_RDM---").Err()
		if err != nil {
			resp.Msg = err.Error()
			return
		}
		removed = append(removed, index)
	} else {
		// replace index value
		err = rdb.LSet(ctx, key, index, value).Err()
		if err != nil {
			resp.Msg = err.Error()
			return
		}
		updated[index] = value
	}

	resp.Success = true
	resp.Data = map[string]any{
		"removed": removed,
		"updated": updated,
	}
	return
}

// SetSetItem add members to set or remove from set
func (c *connectionService) SetSetItem(connName string, db int, key string, remove bool, members []any) (resp types.JSResp) {
	rdb, ctx, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	if remove {
		_, err = rdb.SRem(ctx, key, members...).Result()
	} else {
		_, err = rdb.SAdd(ctx, key, members...).Result()
	}
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	resp.Success = true
	return
}

// UpdateSetItem replace member of set
func (c *connectionService) UpdateSetItem(connName string, db int, key, value, newValue string) (resp types.JSResp) {
	rdb, ctx, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	_, _ = rdb.SRem(ctx, key, value).Result()
	_, err = rdb.SAdd(ctx, key, newValue).Result()
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	resp.Success = true
	return
}

// UpdateZSetValue update value of sorted set member
func (c *connectionService) UpdateZSetValue(connName string, db int, key, value, newValue string, score float64) (resp types.JSResp) {
	rdb, ctx, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	updated := map[string]any{}
	var removed []string
	if len(newValue) <= 0 {
		// blank new value, delete value
		_, err = rdb.ZRem(ctx, key, value).Result()
		if err == nil {
			removed = append(removed, value)
		}
	} else if newValue == value {
		// update score only
		_, err = rdb.ZAdd(ctx, key, redis.Z{
			Score:  score,
			Member: value,
		}).Result()
	} else {
		// remove old value and add new one
		_, err = rdb.ZRem(ctx, key, value).Result()
		if err == nil {
			removed = append(removed, value)
		}

		_, err = rdb.ZAdd(ctx, key, redis.Z{
			Score:  score,
			Member: newValue,
		}).Result()
		if err == nil {
			updated[newValue] = score
		}
	}
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	resp.Success = true
	resp.Data = map[string]any{
		"updated": updated,
		"removed": removed,
	}
	return
}

// AddZSetValue add item to sorted set
func (c *connectionService) AddZSetValue(connName string, db int, key string, action int, valueScore map[string]float64) (resp types.JSResp) {
	rdb, ctx, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	members := maputil.ToSlice(valueScore, func(k string) redis.Z {
		return redis.Z{
			Score:  valueScore[k],
			Member: k,
		}
	})

	switch action {
	case 1:
		// ignore duplicated fields
		_, err = rdb.ZAddNX(ctx, key, members...).Result()
	default:
		// overwrite duplicated fields
		_, err = rdb.ZAdd(ctx, key, members...).Result()
	}
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	resp.Success = true
	return
}

// AddStreamValue add stream field
func (c *connectionService) AddStreamValue(connName string, db int, key, ID string, fieldItems []any) (resp types.JSResp) {
	rdb, ctx, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	_, err = rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: key,
		ID:     ID,
		Values: fieldItems,
	}).Result()
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	resp.Success = true
	return
}

// RemoveStreamValues remove stream values by id
func (c *connectionService) RemoveStreamValues(connName string, db int, key string, IDs []string) (resp types.JSResp) {
	rdb, ctx, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	_, err = rdb.XDel(ctx, key, IDs...).Result()
	resp.Success = true
	return
}

// SetKeyTTL set ttl of key
func (c *connectionService) SetKeyTTL(connName string, db int, key string, ttl int64) (resp types.JSResp) {
	rdb, ctx, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	var expiration time.Duration
	if ttl < 0 {
		if err = rdb.Persist(ctx, key).Err(); err != nil {
			resp.Msg = err.Error()
			return
		}
	} else {
		expiration = time.Duration(ttl) * time.Second
		if err = rdb.Expire(ctx, key, expiration).Err(); err != nil {
			resp.Msg = err.Error()
			return
		}
	}

	resp.Success = true
	return
}

// DeleteKey remove redis key
func (c *connectionService) DeleteKey(connName string, db int, key string) (resp types.JSResp) {
	rdb, ctx, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	var deletedKeys []string

	if strings.HasSuffix(key, "*") {
		// delete by prefix
		var cursor uint64
		for {
			var loadedKey []string
			if loadedKey, cursor, err = rdb.Scan(ctx, cursor, key, 10000).Result(); err != nil {
				resp.Msg = err.Error()
				return
			} else {
				if err = rdb.Del(ctx, loadedKey...).Err(); err != nil {
					resp.Msg = err.Error()
					return
				} else {
					deletedKeys = append(deletedKeys, loadedKey...)
				}
			}

			// no more loadedKey
			if cursor == 0 {
				break
			}
		}
	} else {
		// delete key only
		_, err = rdb.Del(ctx, key).Result()
		if err != nil {
			resp.Msg = err.Error()
			return
		}
		deletedKeys = append(deletedKeys, key)
	}

	resp.Success = true
	resp.Data = map[string]any{
		"deleted": deletedKeys,
	}
	return
}

// RenameKey rename key
func (c *connectionService) RenameKey(connName string, db int, key, newKey string) (resp types.JSResp) {
	rdb, ctx, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	_, err = rdb.RenameNX(ctx, key, newKey).Result()
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	resp.Success = true
	return
}

// GetCmdHistory get redis command history
func (c *connectionService) GetCmdHistory(pageNo, pageSize int) (resp types.JSResp) {
	resp.Success = true
	if pageSize <= 0 || pageNo <= 0 {
		// return all history
		resp.Data = map[string]any{
			"list":     c.cmdHistory,
			"pageNo":   1,
			"pageSize": -1,
		}
	} else {
		total := len(c.cmdHistory)
		startIndex := total / pageSize * (pageNo - 1)
		endIndex := min(startIndex+pageSize, total)
		resp.Data = map[string]any{
			"list":     c.cmdHistory[startIndex:endIndex],
			"pageNo":   pageNo,
			"pageSize": pageSize,
		}
	}
	return
}

// CleanCmdHistory clean redis command history
func (c *connectionService) CleanCmdHistory() (resp types.JSResp) {
	c.cmdHistory = []cmdHistoryItem{}
	resp.Success = true
	return
}

// update or insert key info to database
//func (c *connectionService) updateDBKey(connName string, db int, keys []string, separator string) {
//	dbStruct := map[string]any{}
//	for _, key := range keys {
//		keyPart := strings.Split(key, separator)
//		prefixLen := len(keyPart)-1
//		if prefixLen > 0 {
//			for i := 0; i < prefixLen; i++ {
//				if dbStruct[keyPart[i]]
//				keyPart[i]
//			}
//		}
//		log.Println("key", key)
//	}
//}
