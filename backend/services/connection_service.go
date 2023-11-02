package services

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/ssh"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"tinyrdm/backend/consts"
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
	client     redis.UniversalClient
	ctx        context.Context
	cancelFunc context.CancelFunc
	cursor     map[int]uint64 // current cursor of databases
	stepSize   int64
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

func (c *connectionService) Stop() {
	for _, item := range c.connMap {
		if item.client != nil {
			item.cancelFunc()
			item.client.Close()
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

	var tlsConfig *tls.Config
	if config.SSL.Enable {
		// setup tls config
		var certs []tls.Certificate
		if len(config.SSL.CertFile) > 0 && len(config.SSL.KeyFile) > 0 {
			if cert, err := tls.LoadX509KeyPair(config.SSL.CertFile, config.SSL.KeyFile); err != nil {
				return nil, err
			} else {
				certs = []tls.Certificate{cert}
			}
		}

		var caCertPool *x509.CertPool
		if len(config.SSL.CAFile) > 0 {
			ca, err := os.ReadFile(config.SSL.CAFile)
			if err != nil {
				return nil, err
			}
			caCertPool = x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(ca)
		}

		if len(certs) <= 0 {
			return nil, errors.New("tls config error")
		}

		tlsConfig = &tls.Config{
			RootCAs:            caCertPool,
			InsecureSkipVerify: false,
			Certificates:       certs,
		}
	}

	option := &redis.Options{
		ClientName:   url.QueryEscape(config.Name),
		Addr:         fmt.Sprintf("%s:%d", config.Addr, config.Port),
		Username:     config.Username,
		Password:     config.Password,
		DialTimeout:  time.Duration(config.ConnTimeout) * time.Second,
		ReadTimeout:  time.Duration(config.ExecTimeout) * time.Second,
		WriteTimeout: time.Duration(config.ExecTimeout) * time.Second,
		TLSConfig:    tlsConfig,
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

func (c *connectionService) createRedisClient(config types.ConnectionConfig) (redis.UniversalClient, error) {
	option, err := c.buildOption(config)
	if err != nil {
		return nil, err
	}

	if config.Sentinel.Enable {
		// get master address via sentinel node
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
	if config.Cluster.Enable {
		// connect to cluster
		var slots []redis.ClusterSlot
		if slots, err = rdb.ClusterSlots(c.ctx).Result(); err == nil {
			clusterOptions := &redis.ClusterOptions{
				//NewClient:             nil,
				//MaxRedirects:          0,
				//RouteByLatency:        false,
				//RouteRandomly:         false,
				//ClusterSlots:          nil,
				ClientName:            url.QueryEscape(option.ClientName),
				Dialer:                option.Dialer,
				OnConnect:             option.OnConnect,
				Protocol:              option.Protocol,
				Username:              option.Username,
				Password:              option.Password,
				MaxRetries:            option.MaxRetries,
				MinRetryBackoff:       option.MinRetryBackoff,
				MaxRetryBackoff:       option.MaxRetryBackoff,
				DialTimeout:           option.DialTimeout,
				ContextTimeoutEnabled: option.ContextTimeoutEnabled,
				PoolFIFO:              option.PoolFIFO,
				PoolSize:              option.PoolSize,
				PoolTimeout:           option.PoolTimeout,
				MinIdleConns:          option.MinIdleConns,
				MaxIdleConns:          option.MaxIdleConns,
				ConnMaxIdleTime:       option.ConnMaxIdleTime,
				ConnMaxLifetime:       option.ConnMaxLifetime,
				TLSConfig:             option.TLSConfig,
				DisableIndentity:      option.DisableIndentity,
			}
			if option.Dialer != nil {
				clusterOptions.Dialer = option.Dialer
				clusterOptions.ReadTimeout = -2
				clusterOptions.WriteTimeout = -2
			}
			var addrs []string
			for _, slot := range slots {
				for _, node := range slot.Nodes {
					addrs = append(addrs, node.Addr)
				}
			}
			clusterOptions.Addrs = addrs
			clusterClient := redis.NewClusterClient(clusterOptions)
			return clusterClient, nil
		} else {
			return nil, err
		}
	}

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
	client, err := c.createRedisClient(config)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	defer client.Close()

	if _, err = client.Ping(c.ctx).Result(); err != nil && err != redis.Nil {
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

func (c *connectionService) getConnection(name string) *types.Connection {
	return c.conns.GetConnection(name)
}

// GetConnection get connection profile by name
func (c *connectionService) GetConnection(name string) (resp types.JSResp) {
	conn := c.getConnection(name)
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
	item, err := c.getRedisClient(name, 0)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	// get connection config
	selConn := c.conns.GetConnection(name)

	var totaldb int
	if selConn.DBFilterType == "" || selConn.DBFilterType == "none" {
		// get total databases
		if config, err := client.ConfigGet(ctx, "databases").Result(); err == nil {
			if total, err := strconv.Atoi(config["databases"]); err == nil {
				totaldb = total
			}
		}
	}

	// parse all db, response content like below
	var dbs []types.ConnectionDB
	var clusterKeyCount int64
	cluster, isCluster := client.(*redis.ClusterClient)
	if isCluster {
		var keyCount atomic.Int64
		err = cluster.ForEachMaster(ctx, func(ctx context.Context, cli *redis.Client) error {
			if size, serr := cli.DBSize(ctx).Result(); serr != nil {
				return serr
			} else {
				keyCount.Add(size)
			}
			return nil
		})
		if err != nil {
			resp.Msg = "get db size error:" + err.Error()
			return
		}
		clusterKeyCount = keyCount.Load()

		// only one database in cluster mode
		dbs = []types.ConnectionDB{
			{
				Name:  "db0",
				Index: 0,
				Keys:  int(clusterKeyCount),
			},
		}
	} else {
		// get database info
		var res string
		res, err = client.Info(ctx, "keyspace").Result()
		if err != nil {
			resp.Msg = "get server info fail:" + err.Error()
			return
		}
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
	}

	resp.Success = true
	resp.Data = map[string]any{
		"db":   dbs,
		"view": selConn.KeyView,
	}
	return
}

// CloseConnection close redis server connection
func (c *connectionService) CloseConnection(name string) (resp types.JSResp) {
	item, ok := c.connMap[name]
	if ok {
		delete(c.connMap, name)
		if item.client != nil {
			item.cancelFunc()
			item.client.Close()
		}
	}
	resp.Success = true
	return
}

// get a redis client from local cache or create a new open
// if db >= 0, will also switch to db index
func (c *connectionService) getRedisClient(connName string, db int) (item connectionItem, err error) {
	var ok bool
	var client redis.UniversalClient
	if item, ok = c.connMap[connName]; ok {
		client = item.client
	} else {
		selConn := c.conns.GetConnection(connName)
		if selConn == nil {
			err = fmt.Errorf("no match connection \"%s\"", connName)
			return
		}

		hook := redis2.NewHook(connName, func(cmd string, cost int64) {
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
		})

		client, err = c.createRedisClient(selConn.ConnectionConfig)
		if err != nil {
			err = fmt.Errorf("create conenction error: %s", err.Error())
			return
		}
		// add hook to each node in cluster mode
		var cluster *redis.ClusterClient
		if cluster, ok = client.(*redis.ClusterClient); ok {
			err = cluster.ForEachShard(c.ctx, func(ctx context.Context, cli *redis.Client) error {
				cli.AddHook(hook)
				return nil
			})
			if err != nil {
				err = fmt.Errorf("get cluster nodes error: %s", err.Error())
				return
			}
		} else {
			client.AddHook(hook)
		}

		if _, err = client.Ping(c.ctx).Result(); err != nil && err != redis.Nil {
			err = errors.New("can not connect to redis server:" + err.Error())
			return
		}
		ctx, cancelFunc := context.WithCancel(c.ctx)
		item = connectionItem{
			client:     client,
			ctx:        ctx,
			cancelFunc: cancelFunc,
			cursor:     map[int]uint64{},
			stepSize:   int64(selConn.LoadSize),
		}
		if item.stepSize <= 0 {
			item.stepSize = consts.DEFAULT_LOAD_SIZE
		}
		c.connMap[connName] = item
	}

	if db >= 0 {
		var rdb *redis.Client
		if rdb, ok = client.(*redis.Client); ok && rdb != nil {
			if err = rdb.Do(item.ctx, "select", strconv.Itoa(db)).Err(); err != nil {
				return
			}
		}
	}
	return
}

// save current scan cursor
func (c *connectionService) setClientCursor(connName string, db int, cursor uint64) {
	if _, ok := c.connMap[connName]; ok {
		if cursor == 0 {
			delete(c.connMap[connName].cursor, db)
		} else {
			c.connMap[connName].cursor[db] = cursor
		}
	}
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
	item, err := c.getRedisClient(name, 0)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	// get database info
	res, err := client.Info(ctx).Result()
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
	c.setClientCursor(connName, db, 0)
	return c.LoadNextKeys(connName, db, match, keyType)
}

// scan keys
// @return loaded keys
// @return next cursor
// @return scan error
func (c *connectionService) scanKeys(ctx context.Context, client redis.UniversalClient, match, keyType string, cursor uint64, count int64) ([]any, uint64, error) {
	var err error
	filterType := len(keyType) > 0
	scanSize := int64(Preferences().GetScanSize())
	// define sub scan function
	scan := func(ctx context.Context, cli redis.UniversalClient, appendFunc func(k []any)) error {
		var loadedKey []string
		var scanCount int64
		for {
			if filterType {
				loadedKey, cursor, err = cli.ScanType(ctx, cursor, match, scanSize, keyType).Result()
			} else {
				loadedKey, cursor, err = cli.Scan(ctx, cursor, match, scanSize).Result()
			}
			if err != nil {
				return err
			} else {
				ks := sliceutil.Map(loadedKey, func(i int) any {
					return strutil.EncodeRedisKey(loadedKey[i])
				})
				scanCount += int64(len(ks))
				appendFunc(ks)
			}

			if (count > 0 && scanCount > count) || cursor == 0 {
				break
			}
		}
		return nil
	}

	var keys []any
	if cluster, ok := client.(*redis.ClusterClient); ok {
		// cluster mode
		var mutex sync.Mutex
		err = cluster.ForEachMaster(ctx, func(ctx context.Context, cli *redis.Client) error {
			return scan(ctx, cli, func(k []any) {
				mutex.Lock()
				keys = append(keys, k...)
				mutex.Unlock()
			})
		})
	} else {
		err = scan(ctx, client, func(k []any) {
			keys = append(keys, k...)
		})
	}
	if err != nil {
		return nil, cursor, err
	}
	return keys, cursor, nil
}

// LoadNextKeys load next key from saved cursor
func (c *connectionService) LoadNextKeys(connName string, db int, match, keyType string) (resp types.JSResp) {
	item, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx, count := item.client, item.ctx, item.stepSize
	cursor := item.cursor[db]
	keys, cursor, err := c.scanKeys(ctx, client, match, keyType, cursor, count)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	c.setClientCursor(connName, db, cursor)

	resp.Success = true
	resp.Data = map[string]any{
		"keys": keys,
		"end":  cursor == 0,
	}
	return
}

// LoadAllKeys load all keys
func (c *connectionService) LoadAllKeys(connName string, db int, match, keyType string) (resp types.JSResp) {
	item, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	cursor := item.cursor[db]
	keys, _, err := c.scanKeys(ctx, client, match, keyType, cursor, 0)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	c.setClientCursor(connName, db, 0)

	resp.Success = true
	resp.Data = map[string]any{
		"keys": keys,
	}
	return
}

// GetKeyValue get value by key
func (c *connectionService) GetKeyValue(connName string, db int, k any, viewAs, decodeType string) (resp types.JSResp) {
	item, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(k)
	var keyType string
	var dur time.Duration
	keyType, err = client.Type(ctx, key).Result()
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	if keyType == "none" {
		resp.Msg = "key not exists"
		return
	}

	var ttl int64
	if dur, err = client.TTL(ctx, key).Result(); err != nil {
		ttl = -1
	} else {
		if dur < 0 {
			ttl = -1
		} else {
			ttl = int64(dur.Seconds())
		}
	}

	var value any
	var size, length int64
	var cursor uint64
	switch strings.ToLower(keyType) {
	case "string":
		var str string
		str, err = client.Get(ctx, key).Result()
		value, decodeType, viewAs = strutil.ConvertTo(str, decodeType, viewAs)
		length, _ = client.StrLen(ctx, key).Result()
		size, _ = client.MemoryUsage(ctx, key, 0).Result()
	case "list":
		value, err = client.LRange(ctx, key, 0, -1).Result()
		length, _ = client.LLen(ctx, key).Result()
		size, _ = client.MemoryUsage(ctx, key, 0).Result()
	case "hash":
		//value, err = client.HGetAll(ctx, key).Result()
		items := map[string]string{}
		scanSize := int64(Preferences().GetScanSize())
		for {
			var loadedVal []string
			loadedVal, cursor, err = client.HScan(ctx, key, cursor, "*", scanSize).Result()
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
		length, _ = client.HLen(ctx, key).Result()
		size, _ = client.MemoryUsage(ctx, key, 0).Result()
	case "set":
		//value, err = client.SMembers(ctx, key).Result()
		items := []string{}
		scanSize := int64(Preferences().GetScanSize())
		for {
			var loadedKey []string
			loadedKey, cursor, err = client.SScan(ctx, key, cursor, "*", scanSize).Result()
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
		length, _ = client.SCard(ctx, key).Result()
		size, _ = client.MemoryUsage(ctx, key, 0).Result()
	case "zset":
		//value, err = client.ZRangeWithScores(ctx, key, 0, -1).Result()
		var items []types.ZSetItem
		scanSize := int64(Preferences().GetScanSize())
		for {
			var loadedVal []string
			loadedVal, cursor, err = client.ZScan(ctx, key, cursor, "*", scanSize).Result()
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
		length, _ = client.ZCard(ctx, key).Result()
		size, _ = client.MemoryUsage(ctx, key, 0).Result()
	case "stream":
		var msgs []redis.XMessage
		items := []types.StreamItem{}
		msgs, err = client.XRevRange(ctx, key, "+", "-").Result()
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
		length, _ = client.XLen(ctx, key).Result()
		size, _ = client.MemoryUsage(ctx, key, 0).Result()
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
		"length": length,
		"viewAs": viewAs,
		"decode": decodeType,
	}
	return
}

// SetKeyValue set value by key
// @param ttl <= 0 means keep current ttl
func (c *connectionService) SetKeyValue(connName string, db int, k any, keyType string, value any, ttl int64, viewAs, decode string) (resp types.JSResp) {
	item, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(k)
	var expiration time.Duration
	if ttl < 0 {
		if expiration, err = client.PTTL(ctx, key).Result(); err != nil {
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
			if saveStr, err = strutil.SaveAs(str, viewAs, decode); err != nil {
				resp.Msg = fmt.Sprintf(`save to "%s" type fail: %s`, viewAs, err.Error())
				return
			}
			_, err = client.Set(ctx, key, saveStr, 0).Result()
			// set expiration lonely, not "keepttl"
			if err == nil && expiration > 0 {
				client.Expire(ctx, key, expiration)
			}
		}
	case "list":
		if strs, ok := value.([]any); !ok {
			resp.Msg = "invalid list value"
			return
		} else {
			err = client.LPush(ctx, key, strs...).Err()
			if err == nil && expiration > 0 {
				client.Expire(ctx, key, expiration)
			}
		}
	case "hash":
		if strs, ok := value.([]any); !ok {
			resp.Msg = "invalid hash value"
			return
		} else {
			total := len(strs)
			if total > 1 {
				_, err = client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
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
				err = client.SAdd(ctx, key, strs...).Err()
				if err == nil && expiration > 0 {
					client.Expire(ctx, key, expiration)
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
				err = client.ZAdd(ctx, key, members...).Err()
				if err == nil && expiration > 0 {
					client.Expire(ctx, key, expiration)
				}
			}
		}
	case "stream":
		if strs, ok := value.([]any); !ok {
			resp.Msg = "invalid stream value"
			return
		} else {
			if len(strs) > 2 {
				err = client.XAdd(ctx, &redis.XAddArgs{
					Stream: key,
					ID:     strs[0].(string),
					Values: strs[1:],
				}).Err()
				if err == nil && expiration > 0 {
					client.Expire(ctx, key, expiration)
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
func (c *connectionService) SetHashValue(connName string, db int, k any, field, newField, value string) (resp types.JSResp) {
	item, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(k)
	var removedField []string
	updatedField := map[string]string{}
	if len(field) <= 0 {
		// old filed is empty, add new field
		_, err = client.HSet(ctx, key, newField, value).Result()
		updatedField[newField] = value
	} else if len(newField) <= 0 {
		// new field is empty, delete old field
		_, err = client.HDel(ctx, key, field, value).Result()
		removedField = append(removedField, field)
	} else if field == newField {
		// replace field
		_, err = client.HSet(ctx, key, newField, value).Result()
		updatedField[newField] = value
	} else {
		// remove old field and add new field
		if _, err = client.HDel(ctx, key, field).Result(); err != nil {
			resp.Msg = err.Error()
			return
		}
		_, err = client.HSet(ctx, key, newField, value).Result()
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
func (c *connectionService) AddHashField(connName string, db int, k any, action int, fieldItems []any) (resp types.JSResp) {
	item, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(k)
	updated := map[string]any{}
	switch action {
	case 1:
		// ignore duplicated fields
		for i := 0; i < len(fieldItems); i += 2 {
			_, err = client.HSetNX(ctx, key, fieldItems[i].(string), fieldItems[i+1]).Result()
			if err == nil {
				updated[fieldItems[i].(string)] = fieldItems[i+1]
			}
		}
	default:
		// overwrite duplicated fields
		total := len(fieldItems)
		if total > 1 {
			_, err = client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
				for i := 0; i < total; i += 2 {
					client.HSet(ctx, key, fieldItems[i], fieldItems[i+1])
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
func (c *connectionService) AddListItem(connName string, db int, k any, action int, items []any) (resp types.JSResp) {
	item, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(k)
	var leftPush, rightPush []any
	switch action {
	case 0:
		// push to head
		_, err = client.LPush(ctx, key, items...).Result()
		leftPush = append(leftPush, items...)
	default:
		// append to tail
		_, err = client.RPush(ctx, key, items...).Result()
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
func (c *connectionService) SetListItem(connName string, db int, k any, index int64, value string) (resp types.JSResp) {
	item, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(k)
	var removed []int64
	updated := map[int64]string{}
	if len(value) <= 0 {
		// remove from list
		err = client.LSet(ctx, key, index, "---VALUE_REMOVED_BY_TINY_RDM---").Err()
		if err != nil {
			resp.Msg = err.Error()
			return
		}

		err = client.LRem(ctx, key, 1, "---VALUE_REMOVED_BY_TINY_RDM---").Err()
		if err != nil {
			resp.Msg = err.Error()
			return
		}
		removed = append(removed, index)
	} else {
		// replace index value
		err = client.LSet(ctx, key, index, value).Err()
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
func (c *connectionService) SetSetItem(connName string, db int, k any, remove bool, members []any) (resp types.JSResp) {
	item, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(k)
	if remove {
		_, err = client.SRem(ctx, key, members...).Result()
	} else {
		_, err = client.SAdd(ctx, key, members...).Result()
	}
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	resp.Success = true
	return
}

// UpdateSetItem replace member of set
func (c *connectionService) UpdateSetItem(connName string, db int, k any, value, newValue string) (resp types.JSResp) {
	item, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(k)
	_, _ = client.SRem(ctx, key, value).Result()
	_, err = client.SAdd(ctx, key, newValue).Result()
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	resp.Success = true
	return
}

// UpdateZSetValue update value of sorted set member
func (c *connectionService) UpdateZSetValue(connName string, db int, k any, value, newValue string, score float64) (resp types.JSResp) {
	item, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(k)
	updated := map[string]any{}
	var removed []string
	if len(newValue) <= 0 {
		// blank new value, delete value
		_, err = client.ZRem(ctx, key, value).Result()
		if err == nil {
			removed = append(removed, value)
		}
	} else if newValue == value {
		// update score only
		_, err = client.ZAdd(ctx, key, redis.Z{
			Score:  score,
			Member: value,
		}).Result()
	} else {
		// remove old value and add new one
		_, err = client.ZRem(ctx, key, value).Result()
		if err == nil {
			removed = append(removed, value)
		}

		_, err = client.ZAdd(ctx, key, redis.Z{
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
func (c *connectionService) AddZSetValue(connName string, db int, k any, action int, valueScore map[string]float64) (resp types.JSResp) {
	item, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(k)
	members := maputil.ToSlice(valueScore, func(k string) redis.Z {
		return redis.Z{
			Score:  valueScore[k],
			Member: k,
		}
	})

	switch action {
	case 1:
		// ignore duplicated fields
		_, err = client.ZAddNX(ctx, key, members...).Result()
	default:
		// overwrite duplicated fields
		_, err = client.ZAdd(ctx, key, members...).Result()
	}
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	resp.Success = true
	return
}

// AddStreamValue add stream field
func (c *connectionService) AddStreamValue(connName string, db int, k any, ID string, fieldItems []any) (resp types.JSResp) {
	item, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(k)
	_, err = client.XAdd(ctx, &redis.XAddArgs{
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
func (c *connectionService) RemoveStreamValues(connName string, db int, k any, IDs []string) (resp types.JSResp) {
	item, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(k)
	_, err = client.XDel(ctx, key, IDs...).Result()
	resp.Success = true
	return
}

// SetKeyTTL set ttl of key
func (c *connectionService) SetKeyTTL(connName string, db int, k any, ttl int64) (resp types.JSResp) {
	item, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(k)
	var expiration time.Duration
	if ttl < 0 {
		if err = client.Persist(ctx, key).Err(); err != nil {
			resp.Msg = err.Error()
			return
		}
	} else {
		expiration = time.Duration(ttl) * time.Second
		if err = client.Expire(ctx, key, expiration).Err(); err != nil {
			resp.Msg = err.Error()
			return
		}
	}

	resp.Success = true
	return
}

// DeleteKey remove redis key
func (c *connectionService) DeleteKey(connName string, db int, k any, async bool) (resp types.JSResp) {
	item, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(k)
	var deletedKeys []string
	if strings.HasSuffix(key, "*") {
		// delete by prefix
		var mutex sync.Mutex
		del := func(ctx context.Context, cli redis.UniversalClient) error {
			handleDel := func(ks []string) error {
				pipe := cli.Pipeline()
				for _, k2 := range ks {
					if async {
						cli.Unlink(ctx, k2)
					} else {
						cli.Del(ctx, k2)
					}
				}
				pipe.Exec(ctx)

				mutex.Lock()
				deletedKeys = append(deletedKeys, ks...)
				mutex.Unlock()

				return nil
			}

			scanSize := int64(Preferences().GetScanSize())
			iter := cli.Scan(ctx, 0, key, scanSize).Iterator()
			resultKeys := make([]string, 0, 100)
			for iter.Next(ctx) {
				resultKeys = append(resultKeys, iter.Val())
				if len(resultKeys) >= 3 {
					handleDel(resultKeys)
					resultKeys = resultKeys[:0:cap(resultKeys)]
				}
			}

			if len(resultKeys) > 0 {
				handleDel(resultKeys)
			}
			return nil
		}

		if cluster, ok := client.(*redis.ClusterClient); ok {
			// cluster mode
			err = cluster.ForEachMaster(ctx, func(ctx context.Context, cli *redis.Client) error {
				return del(ctx, cli)
			})
		} else {
			err = del(ctx, client)
		}

		if err != nil {
			resp.Msg = err.Error()
			return
		}
	} else {
		// delete key only
		if async {
			if _, err = client.Unlink(ctx, key).Result(); err != nil {
				resp.Msg = err.Error()
				return
			}
		} else {
			if _, err = client.Del(ctx, key).Result(); err != nil {
				resp.Msg = err.Error()
				return
			}
		}
		deletedKeys = append(deletedKeys, key)
	}

	resp.Success = true
	resp.Data = map[string]any{
		"deleted": deletedKeys,
	}
	return
}

// FlushDB flush database
func (c *connectionService) FlushDB(connName string, db int, async bool) (resp types.JSResp) {
	item, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	flush := func(ctx context.Context, cli redis.UniversalClient) {
		cli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Select(ctx, db)
			if async {
				pipe.FlushDBAsync(ctx)
			} else {
				pipe.FlushDB(ctx)
			}
			return nil
		})
	}

	client, ctx := item.client, item.ctx
	if cluster, ok := client.(*redis.ClusterClient); ok {
		// cluster mode
		err = cluster.ForEachMaster(ctx, func(ctx context.Context, cli *redis.Client) error {
			flush(ctx, cli)
			return nil
		})
	} else {
		flush(ctx, client)
	}

	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Success = true
	return
}

// RenameKey rename key
func (c *connectionService) RenameKey(connName string, db int, key, newKey string) (resp types.JSResp) {
	item, err := c.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	if _, ok := client.(*redis.ClusterClient); ok {
		resp.Msg = "RENAME not support in cluster mode yet"
		return
	}

	if _, err = client.RenameNX(ctx, key, newKey).Result(); err != nil {
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
