package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"tinyrdm/backend/consts"
	"tinyrdm/backend/types"
	"tinyrdm/backend/utils/coll"
	maputil "tinyrdm/backend/utils/map"
	redis2 "tinyrdm/backend/utils/redis"
	sliceutil "tinyrdm/backend/utils/slice"
	strutil "tinyrdm/backend/utils/string"
)

type slowLogItem struct {
	Timestamp int64  `json:"timestamp"`
	Client    string `json:"client"`
	Addr      string `json:"addr"`
	Cmd       string `json:"cmd"`
	Cost      int64  `json:"cost"`
}

type entryCursor struct {
	DB      int
	Type    string
	Key     string
	Pattern string
	Cursor  uint64
	XLast   string // last stream pos
}

type connectionItem struct {
	client      redis.UniversalClient
	ctx         context.Context
	cancelFunc  context.CancelFunc
	cursor      map[int]uint64      // current cursor of databases
	entryCursor map[int]entryCursor // current entry cursor of databases
	stepSize    int64
}

type browserService struct {
	ctx        context.Context
	connMap    map[string]connectionItem
	cmdHistory []cmdHistoryItem
}

var browser *browserService
var onceBrowser sync.Once

func Browser() *browserService {
	if browser == nil {
		onceBrowser.Do(func() {
			browser = &browserService{
				connMap: map[string]connectionItem{},
			}
		})
	}
	return browser
}

func (b *browserService) Start(ctx context.Context) {
	b.ctx = ctx
}

func (b *browserService) Stop() {
	for _, item := range b.connMap {
		if item.client != nil {
			item.cancelFunc()
			item.client.Close()
		}
	}
	b.connMap = map[string]connectionItem{}
}

// OpenConnection open redis server connection
func (b *browserService) OpenConnection(name string) (resp types.JSResp) {
	item, err := b.getRedisClient(name, 0)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	// get connection config
	selConn := Connection().getConnection(name)

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
		info := b.parseInfo(res)

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
				dbInfo := b.parseDBItemInfo(dbInfoStr)
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
func (b *browserService) CloseConnection(name string) (resp types.JSResp) {
	item, ok := b.connMap[name]
	if ok {
		delete(b.connMap, name)
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
func (b *browserService) getRedisClient(connName string, db int) (item connectionItem, err error) {
	var ok bool
	var client redis.UniversalClient
	if item, ok = b.connMap[connName]; ok {
		client = item.client
	} else {
		selConn := Connection().getConnection(connName)
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
			b.cmdHistory = append(b.cmdHistory, cmdHistoryItem{
				Timestamp: now.UnixMilli(),
				Server:    connName,
				Cmd:       cmd,
				Cost:      cost,
			})
		})

		client, err = Connection().createRedisClient(selConn.ConnectionConfig)
		if err != nil {
			err = fmt.Errorf("create conenction error: %s", err.Error())
			return
		}
		// add hook to each node in cluster mode
		var cluster *redis.ClusterClient
		if cluster, ok = client.(*redis.ClusterClient); ok {
			err = cluster.ForEachShard(b.ctx, func(ctx context.Context, cli *redis.Client) error {
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

		if _, err = client.Ping(b.ctx).Result(); err != nil && err != redis.Nil {
			err = errors.New("can not connect to redis server:" + err.Error())
			return
		}
		ctx, cancelFunc := context.WithCancel(b.ctx)
		item = connectionItem{
			client:      client,
			ctx:         ctx,
			cancelFunc:  cancelFunc,
			cursor:      map[int]uint64{},
			entryCursor: map[int]entryCursor{},
			stepSize:    int64(selConn.LoadSize),
		}
		if item.stepSize <= 0 {
			item.stepSize = consts.DEFAULT_LOAD_SIZE
		}
		b.connMap[connName] = item
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

// load current database size
func (b *browserService) loadDBSize(ctx context.Context, client redis.UniversalClient) int64 {
	if cluster, isCluster := client.(*redis.ClusterClient); isCluster {
		var keyCount atomic.Int64
		cluster.ForEachMaster(ctx, func(ctx context.Context, cli *redis.Client) error {
			if size, serr := cli.DBSize(ctx).Result(); serr != nil {
				return serr
			} else {
				keyCount.Add(size)
			}
			return nil
		})
		return keyCount.Load()
	} else {
		keyCount, _ := client.DBSize(ctx).Result()
		return keyCount
	}
}

// save current scan cursor
func (b *browserService) setClientCursor(connName string, db int, cursor uint64) {
	if _, ok := b.connMap[connName]; ok {
		if cursor == 0 {
			delete(b.connMap[connName].cursor, db)
		} else {
			b.connMap[connName].cursor[db] = cursor
		}
	}
}

// parse command response content which use "redis info"
// # Keyspace\r\ndb0:keys=2,expires=1,avg_ttl=1877111749\r\ndb1:keys=33,expires=0,avg_ttl=0\r\ndb3:keys=17,expires=0,avg_ttl=0\r\ndb5:keys=3,expires=0,avg_ttl=0\r\n
func (b *browserService) parseInfo(info string) map[string]map[string]string {
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
func (b *browserService) parseDBItemInfo(info string) map[string]int {
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
func (b *browserService) ServerInfo(name string) (resp types.JSResp) {
	item, err := b.getRedisClient(name, 0)
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
	resp.Data = b.parseInfo(res)
	return
}

// OpenDatabase open select database, and list all keys
// @param path contain connection name and db name
func (b *browserService) OpenDatabase(connName string, db int, match string, keyType string) (resp types.JSResp) {
	b.setClientCursor(connName, db, 0)
	return b.LoadNextKeys(connName, db, match, keyType)
}

// scan keys
// @return loaded keys
// @return next cursor
// @return scan error
func (b *browserService) scanKeys(ctx context.Context, client redis.UniversalClient, match, keyType string, cursor uint64, count int64) ([]any, uint64, error) {
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
func (b *browserService) LoadNextKeys(connName string, db int, match, keyType string) (resp types.JSResp) {
	item, err := b.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx, count := item.client, item.ctx, item.stepSize
	cursor := item.cursor[db]
	keys, cursor, err := b.scanKeys(ctx, client, match, keyType, cursor, count)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	b.setClientCursor(connName, db, cursor)
	maxKeys := b.loadDBSize(ctx, client)

	resp.Success = true
	resp.Data = map[string]any{
		"keys":    keys,
		"end":     cursor == 0,
		"maxKeys": maxKeys,
	}
	return
}

// LoadAllKeys load all keys
func (b *browserService) LoadAllKeys(connName string, db int, match, keyType string) (resp types.JSResp) {
	item, err := b.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	cursor := item.cursor[db]
	keys, _, err := b.scanKeys(ctx, client, match, keyType, cursor, 0)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	b.setClientCursor(connName, db, 0)
	maxKeys := b.loadDBSize(ctx, client)

	resp.Success = true
	resp.Data = map[string]any{
		"keys":    keys,
		"maxKeys": maxKeys,
	}
	return
}

// GetKeySummary get key summary info
func (b *browserService) GetKeySummary(param types.KeySummaryParam) (resp types.JSResp) {
	item, err := b.getRedisClient(param.Server, param.DB)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(param.Key)
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

	var data types.KeySummary
	data.Type = strings.ToLower(keyType)
	if dur, err = client.TTL(ctx, key).Result(); err != nil {
		data.TTL = -1
	} else {
		if dur < 0 {
			data.TTL = -1
		} else {
			data.TTL = int64(dur.Seconds())
		}
	}

	data.Size, _ = client.MemoryUsage(ctx, key, 0).Result()
	switch data.Type {
	case "string":
		data.Length, _ = client.StrLen(ctx, key).Result()
	case "list":
		data.Length, _ = client.LLen(ctx, key).Result()
	case "hash":
		data.Length, _ = client.HLen(ctx, key).Result()
	case "set":
		data.Length, _ = client.SCard(ctx, key).Result()
	case "zset":
		data.Length, _ = client.ZCard(ctx, key).Result()
	case "stream":
		data.Length, _ = client.XLen(ctx, key).Result()
	default:
		resp.Msg = "unknown key type"
		return
	}

	resp.Success = true
	resp.Data = data
	return
}

// GetKeyDetail get key detail
func (b *browserService) GetKeyDetail(param types.KeyDetailParam) (resp types.JSResp) {
	item, err := b.getRedisClient(param.Server, param.DB)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx, entryCors := item.client, item.ctx, item.entryCursor
	key := strutil.DecodeRedisKey(param.Key)
	var keyType string
	keyType, err = client.Type(ctx, key).Result()
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	if keyType == "none" {
		resp.Msg = "key not exists"
		return
	}
	var doConvert bool
	if (len(param.Decode) > 0 && param.Decode != types.DECODE_NONE) ||
		(len(param.Format) > 0 && param.Format != types.FORMAT_RAW) {
		doConvert = true
	}

	var data types.KeyDetail
	//var cursor uint64
	matchPattern := param.MatchPattern
	if len(matchPattern) <= 0 {
		matchPattern = "*"
	}

	// define get entry cursor function
	getEntryCursor := func() (uint64, string) {
		if entry, ok := entryCors[param.DB]; !ok || entry.Key != key || entry.Pattern != matchPattern {
			// not the same key or match pattern, reset cursor
			entry = entryCursor{
				DB:      param.DB,
				Key:     key,
				Pattern: matchPattern,
				Cursor:  0,
			}
			entryCors[param.DB] = entry
			return 0, ""
		} else {
			return entry.Cursor, entry.XLast
		}
	}
	// define set entry cursor function
	setEntryCursor := func(cursor uint64) {
		entryCors[param.DB] = entryCursor{
			DB:      param.DB,
			Type:    "",
			Key:     key,
			Pattern: matchPattern,
			Cursor:  cursor,
		}
	}
	// define set last stream pos function
	setEntryXLast := func(last string) {
		entryCors[param.DB] = entryCursor{
			DB:      param.DB,
			Type:    "",
			Key:     key,
			Pattern: matchPattern,
			XLast:   last,
		}
	}

	switch strings.ToLower(keyType) {
	case "string":
		var str string
		str, err = client.Get(ctx, key).Result()
		data.Value = strutil.EncodeRedisKey(str)
		//data.Value, data.Decode, data.Format = strutil.ConvertTo(str, param.Decode, param.Format)

	case "list":
		loadListHandle := func() ([]types.ListEntryItem, bool, error) {
			var loadVal []string
			var cursor uint64
			var subErr error
			if param.Full {
				// load all
				cursor = 0
				loadVal, subErr = client.LRange(ctx, key, 0, -1).Result()
			} else {
				cursor, _ = getEntryCursor()
				scanSize := int64(Preferences().GetScanSize())
				loadVal, subErr = client.LRange(ctx, key, int64(cursor), int64(cursor)+scanSize-1).Result()
				cursor = cursor + uint64(scanSize)
				if len(loadVal) < int(scanSize) {
					cursor = 0
				}
			}
			setEntryCursor(cursor)

			items := make([]types.ListEntryItem, len(loadVal))
			for i, val := range loadVal {
				items[i].Value = val
				if doConvert {
					if dv, _, _ := strutil.ConvertTo(val, param.Decode, param.Format); dv != val {
						items[i].DisplayValue = dv
					}
				}
			}
			if subErr != nil {
				return items, false, subErr
			}
			return items, cursor == 0, nil
		}

		data.Value, data.End, err = loadListHandle()
		data.Decode, data.Format = param.Decode, param.Format
		if err != nil {
			resp.Msg = err.Error()
			return
		}

	case "hash":
		loadHashHandle := func() ([]types.HashEntryItem, bool, error) {
			var items []types.HashEntryItem
			var loadedVal []string
			var cursor uint64
			var subErr error
			scanSize := int64(Preferences().GetScanSize())
			if param.Full {
				// load all
				cursor = 0
				for {
					loadedVal, cursor, subErr = client.HScan(ctx, key, cursor, "*", scanSize).Result()
					if subErr != nil {
						return nil, false, subErr
					}
					for i := 0; i < len(loadedVal); i += 2 {
						items = append(items, types.HashEntryItem{
							Key:   loadedVal[i],
							Value: strutil.EncodeRedisKey(loadedVal[i+1]),
						})
						if doConvert {
							if dv, _, _ := strutil.ConvertTo(loadedVal[i+1], param.Decode, param.Format); dv != loadedVal[i+1] {
								items[len(items)-1].DisplayValue = dv
							}
						}
					}
					if cursor == 0 {
						break
					}
				}
			} else {
				cursor, _ = getEntryCursor()
				loadedVal, cursor, subErr = client.HScan(ctx, key, cursor, matchPattern, scanSize).Result()
				if subErr != nil {
					return nil, false, subErr
				}
				loadedLen := len(loadedVal)
				items = make([]types.HashEntryItem, loadedLen/2)
				for i := 0; i < loadedLen; i += 2 {
					items[i/2].Key = loadedVal[i]
					items[i/2].Value = strutil.EncodeRedisKey(loadedVal[i+1])
					if doConvert {
						if dv, _, _ := strutil.ConvertTo(loadedVal[i+1], param.Decode, param.Format); dv != loadedVal[i+1] {
							items[i/2].DisplayValue = dv
						}
					}
				}
			}
			setEntryCursor(cursor)
			return items, cursor == 0, nil
		}

		data.Value, data.End, err = loadHashHandle()
		data.Decode, data.Format = param.Decode, param.Format
		if err != nil {
			resp.Msg = err.Error()
			return
		}

	case "set":
		loadSetHandle := func() ([]types.SetEntryItem, bool, error) {
			var items []types.SetEntryItem
			var cursor uint64
			var subErr error
			var loadedKey []string
			scanSize := int64(Preferences().GetScanSize())
			if param.Full {
				// load all
				cursor = 0
				for {
					loadedKey, cursor, subErr = client.SScan(ctx, key, cursor, param.MatchPattern, scanSize).Result()
					if subErr != nil {
						return items, false, subErr
					}
					for _, val := range loadedKey {
						items = append(items, types.SetEntryItem{
							Value: val,
						})
						if doConvert {
							if dv, _, _ := strutil.ConvertTo(val, param.Decode, param.Format); dv != val {
								items[len(items)-1].DisplayValue = dv
							}
						}
					}
					if cursor == 0 {
						break
					}
				}
			} else {
				cursor, _ = getEntryCursor()
				loadedKey, cursor, subErr = client.SScan(ctx, key, cursor, param.MatchPattern, scanSize).Result()
				items = make([]types.SetEntryItem, len(loadedKey))
				for i, val := range loadedKey {
					items[i].Value = val
					if doConvert {
						if dv, _, _ := strutil.ConvertTo(val, param.Decode, param.Format); dv != val {
							items[i].DisplayValue = dv
						}
					}
				}
			}
			setEntryCursor(cursor)
			return items, cursor == 0, nil
		}

		data.Value, data.End, err = loadSetHandle()
		data.Decode, data.Format = param.Decode, param.Format
		if err != nil {
			resp.Msg = err.Error()
			return
		}

	case "zset":
		loadZSetHandle := func() ([]types.ZSetEntryItem, bool, error) {
			var items []types.ZSetEntryItem
			var cursor uint64
			scanSize := int64(Preferences().GetScanSize())
			var loadedVal []string
			if param.Full {
				// load all
				cursor = 0
				for {
					loadedVal, cursor, err = client.ZScan(ctx, key, cursor, param.MatchPattern, scanSize).Result()
					if err != nil {
						return items, false, err
					}
					var score float64
					for i := 0; i < len(loadedVal); i += 2 {
						if score, err = strconv.ParseFloat(loadedVal[i+1], 64); err == nil {
							items = append(items, types.ZSetEntryItem{
								Value: loadedVal[i],
								Score: score,
							})
							if doConvert {
								if dv, _, _ := strutil.ConvertTo(loadedVal[i], param.Decode, param.Format); dv != loadedVal[i] {
									items[len(items)-1].DisplayValue = dv
								}
							}
						}
					}
					if cursor == 0 {
						break
					}
				}
			} else {
				cursor, _ = getEntryCursor()
				loadedVal, cursor, err = client.ZScan(ctx, key, cursor, param.MatchPattern, scanSize).Result()
				loadedLen := len(loadedVal)
				items = make([]types.ZSetEntryItem, loadedLen/2)
				var score float64
				for i := 0; i < loadedLen; i += 2 {
					if score, err = strconv.ParseFloat(loadedVal[i+1], 64); err == nil {
						items[i/2].Score = score
						items[i/2].Value = loadedVal[i]
						if doConvert {
							if dv, _, _ := strutil.ConvertTo(loadedVal[i], param.Decode, param.Format); dv != loadedVal[i] {
								items[i/2].DisplayValue = dv
							}
						}
					}
				}
			}
			setEntryCursor(cursor)
			return items, cursor == 0, nil
		}

		data.Value, data.End, err = loadZSetHandle()
		data.Decode, data.Format = param.Decode, param.Format
		if err != nil {
			resp.Msg = err.Error()
			return
		}

	case "stream":
		loadStreamHandle := func() ([]types.StreamEntryItem, bool, error) {
			var msgs []redis.XMessage
			var items []types.StreamEntryItem
			var last string
			if param.Full {
				// load all
				last = ""
				msgs, err = client.XRevRange(ctx, key, "+", "-").Result()
			} else {
				scanSize := int64(Preferences().GetScanSize())
				_, last = getEntryCursor()
				if len(last) <= 0 {
					last = "+"
				}
				if last != "+" {
					// add 1 more item when continue scan
					msgs, err = client.XRevRangeN(ctx, key, last, "-", scanSize+1).Result()
					msgs = msgs[1:]
				} else {
					msgs, err = client.XRevRangeN(ctx, key, last, "-", scanSize).Result()
				}
				scanCount := len(msgs)
				if scanCount <= 0 || scanCount < int(scanSize) {
					last = ""
				} else if scanCount > 0 {
					last = msgs[scanCount-1].ID
				}
			}
			setEntryXLast(last)
			for _, msg := range msgs {
				items = append(items, types.StreamEntryItem{
					ID:    msg.ID,
					Value: msg.Values,
				})
			}
			if err != nil {
				return items, false, err
			}
			return items, last == "", nil
		}

		data.Value, data.End, err = loadStreamHandle()
		data.Decode, data.Format = param.Decode, param.Format
		if err != nil {
			resp.Msg = err.Error()
			return
		}
	}
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Success = true
	resp.Data = data
	return
}

// ConvertValue convert value with decode method and format
// blank decode indicate auto decode
// blank format indicate auto format
func (b *browserService) ConvertValue(value any, decode, format string) (resp types.JSResp) {
	str := strutil.DecodeRedisKey(value)
	value, decode, format = strutil.ConvertTo(str, decode, format)
	resp.Success = true
	resp.Data = map[string]any{
		"value":  value,
		"decode": decode,
		"format": format,
	}
	return
}

// SetKeyValue set value by key
// @param ttl <= 0 means keep current ttl
func (b *browserService) SetKeyValue(param types.SetKeyParam) (resp types.JSResp) {
	item, err := b.getRedisClient(param.Server, param.DB)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(param.Key)
	var expiration time.Duration
	if param.TTL < 0 {
		if expiration, err = client.PTTL(ctx, key).Result(); err != nil {
			expiration = redis.KeepTTL
		}
	} else {
		expiration = time.Duration(param.TTL) * time.Second
	}
	// use default decode type and format
	if len(param.Decode) <= 0 {
		param.Decode = types.DECODE_NONE
	}
	if len(param.Format) <= 0 {
		param.Format = types.FORMAT_RAW
	}
	switch strings.ToLower(param.KeyType) {
	case "string":
		if str, ok := param.Value.(string); !ok {
			resp.Msg = "invalid string value"
			return
		} else {
			var saveStr string
			if saveStr, err = strutil.SaveAs(str, param.Format, param.Decode); err != nil {
				resp.Msg = fmt.Sprintf(`save to type "%s" fail: %s`, param.Format, err.Error())
				return
			}
			_, err = client.Set(ctx, key, saveStr, 0).Result()
			// set expiration lonely, not "keepttl"
			if err == nil && expiration > 0 {
				client.Expire(ctx, key, expiration)
			}
		}
	case "list":
		if strs, ok := param.Value.([]any); !ok {
			resp.Msg = "invalid list value"
			return
		} else {
			err = client.LPush(ctx, key, strs...).Err()
			if err == nil && expiration > 0 {
				client.Expire(ctx, key, expiration)
			}
		}
	case "hash":
		if strs, ok := param.Value.([]any); !ok {
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
		if strs, ok := param.Value.([]any); !ok || len(strs) <= 0 {
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
		if strs, ok := param.Value.([]any); !ok || len(strs) <= 0 {
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
		if strs, ok := param.Value.([]any); !ok {
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
		"value": param.Value,
	}
	return
}

// SetHashValue set hash field
func (b *browserService) SetHashValue(param types.SetHashParam) (resp types.JSResp) {
	item, err := b.getRedisClient(param.Server, param.DB)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(param.Key)
	str := strutil.DecodeRedisKey(param.Value)
	var saveStr string
	if saveStr, err = strutil.SaveAs(str, param.Format, param.Decode); err != nil {
		resp.Msg = fmt.Sprintf(`save to type "%s" fail: %s`, param.Format, err.Error())
		return
	}
	var removedField []string
	updatedField := map[string]any{}
	replacedField := map[string]any{}
	if len(param.Field) <= 0 {
		// old filed is empty, add new field
		_, err = client.HSet(ctx, key, param.NewField, saveStr).Result()
		updatedField[param.NewField] = saveStr
	} else if len(param.NewField) <= 0 {
		// new field is empty, delete old field
		_, err = client.HDel(ctx, key, param.Field).Result()
		removedField = append(removedField, param.Field)
	} else if param.Field == param.NewField {
		// update field value
		_, err = client.HSet(ctx, key, param.Field, saveStr).Result()
		updatedField[param.NewField] = saveStr
	} else {
		// remove old field and add new field
		if _, err = client.HDel(ctx, key, param.Field).Result(); err != nil {
			resp.Msg = err.Error()
			return
		}
		_, err = client.HSet(ctx, key, param.NewField, saveStr).Result()
		removedField = append(removedField, param.Field)
		updatedField[param.NewField] = saveStr
		replacedField[param.Field] = param.NewField
	}
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	resp.Success = true
	resp.Data = map[string]any{
		"removed":  removedField,
		"updated":  updatedField,
		"replaced": replacedField,
	}
	return
}

// AddHashField add or update hash field
func (b *browserService) AddHashField(connName string, db int, k any, action int, fieldItems []any) (resp types.JSResp) {
	item, err := b.getRedisClient(connName, db)
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
func (b *browserService) AddListItem(connName string, db int, k any, action int, items []any) (resp types.JSResp) {
	item, err := b.getRedisClient(connName, db)
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
func (b *browserService) SetListItem(param types.SetListParam) (resp types.JSResp) {
	item, err := b.getRedisClient(param.Server, param.DB)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(param.Key)
	str := strutil.DecodeRedisKey(param.Value)
	var removed []int64
	updated := map[int64]string{}
	if len(str) <= 0 {
		// remove from list
		err = client.LSet(ctx, key, param.Index, "---VALUE_REMOVED_BY_TINY_RDM---").Err()
		if err != nil {
			resp.Msg = err.Error()
			return
		}

		err = client.LRem(ctx, key, 1, "---VALUE_REMOVED_BY_TINY_RDM---").Err()
		if err != nil {
			resp.Msg = err.Error()
			return
		}
		removed = append(removed, param.Index)
	} else {
		// replace index value
		var saveStr string
		if saveStr, err = strutil.SaveAs(str, param.Format, param.Decode); err != nil {
			resp.Msg = fmt.Sprintf(`save to type "%s" fail: %s`, param.Format, err.Error())
			return
		}
		err = client.LSet(ctx, key, param.Index, saveStr).Err()
		if err != nil {
			resp.Msg = err.Error()
			return
		}
		updated[param.Index] = saveStr
	}

	resp.Success = true
	resp.Data = map[string]any{
		"removed": removed,
		"updated": updated,
	}
	return
}

// SetSetItem add members to set or remove from set
func (b *browserService) SetSetItem(server string, db int, k any, remove bool, members []any) (resp types.JSResp) {
	item, err := b.getRedisClient(server, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(k)
	var affected int64
	if remove {
		affected, err = client.SRem(ctx, key, members...).Result()
	} else {
		affected, err = client.SAdd(ctx, key, members...).Result()
	}
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	resp.Success = true
	resp.Data = map[string]any{
		"affected": affected,
	}
	return
}

// UpdateSetItem replace member of set
func (b *browserService) UpdateSetItem(param types.SetSetParam) (resp types.JSResp) {
	item, err := b.getRedisClient(param.Server, param.DB)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(param.Key)
	// remove old value
	str := strutil.DecodeRedisKey(param.Value)
	_, _ = client.SRem(ctx, key, str).Result()

	// insert new value
	str = strutil.DecodeRedisKey(param.NewValue)
	var saveStr string
	if saveStr, err = strutil.SaveAs(str, param.Format, param.Decode); err != nil {
		resp.Msg = fmt.Sprintf(`save to type "%s" fail: %s`, param.Format, err.Error())
		return
	}
	_, err = client.SAdd(ctx, key, saveStr).Result()
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	resp.Success = true
	resp.Data = map[string]any{
		"added": saveStr,
	}
	return
}

// UpdateZSetValue update value of sorted set member
func (b *browserService) UpdateZSetValue(param types.SetZSetParam) (resp types.JSResp) {
	item, err := b.getRedisClient(param.Server, param.DB)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(param.Key)
	val, newVal := strutil.DecodeRedisKey(param.Value), strutil.DecodeRedisKey(param.NewValue)
	updated := map[string]float64{}
	var removed []string
	if len(newVal) <= 0 {
		// no new value, delete value
		_, err = client.ZRem(ctx, key, val).Result()
		if err == nil {
			removed = append(removed, val)
		}
	} else {
		var saveVal string
		if saveVal, err = strutil.SaveAs(newVal, param.Format, param.Decode); err != nil {
			resp.Msg = fmt.Sprintf(`save to type "%s" fail: %s`, param.Format, err.Error())
			return
		}

		if saveVal == val {
			// update score only
			_, err = client.ZAdd(ctx, key, redis.Z{
				Score:  param.Score,
				Member: saveVal,
			}).Result()
			if err == nil {
				updated[saveVal] = param.Score
			}
		} else {
			// remove old value and add new one
			_, err = client.ZRem(ctx, key, val).Result()
			if err == nil {
				removed = append(removed, val)
			}

			_, err = client.ZAdd(ctx, key, redis.Z{
				Score:  param.Score,
				Member: saveVal,
			}).Result()
			if err == nil {
				updated[saveVal] = param.Score
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
		"removed": removed,
	}
	return
}

// AddZSetValue add item to sorted set
func (b *browserService) AddZSetValue(connName string, db int, k any, action int, valueScore map[string]float64) (resp types.JSResp) {
	item, err := b.getRedisClient(connName, db)
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
func (b *browserService) AddStreamValue(connName string, db int, k any, ID string, fieldItems []any) (resp types.JSResp) {
	item, err := b.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(k)
	var updateID string
	updateID, err = client.XAdd(ctx, &redis.XAddArgs{
		Stream: key,
		ID:     ID,
		Values: fieldItems,
	}).Result()
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	resp.Success = true
	resp.Data = map[string]any{
		"updateID": updateID,
	}
	return
}

// RemoveStreamValues remove stream values by id
func (b *browserService) RemoveStreamValues(connName string, db int, k any, IDs []string) (resp types.JSResp) {
	item, err := b.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	key := strutil.DecodeRedisKey(k)

	var affected int64
	affected, err = client.XDel(ctx, key, IDs...).Result()
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	resp.Success = true
	resp.Data = map[string]any{
		"affected": affected,
	}
	return
}

// SetKeyTTL set ttl of key
func (b *browserService) SetKeyTTL(connName string, db int, k any, ttl int64) (resp types.JSResp) {
	item, err := b.getRedisClient(connName, db)
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
func (b *browserService) DeleteKey(connName string, db int, k any, async bool) (resp types.JSResp) {
	item, err := b.getRedisClient(connName, db)
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
		supportUnlink := true
		del := func(ctx context.Context, cli redis.UniversalClient) error {
			handleDel := func(ks []string) error {
				var delErr error
				if async && supportUnlink {
					supportUnlink = false
					if delErr = cli.Unlink(ctx, ks...).Err(); delErr != nil {
						// not support unlink? try del command
						delErr = cli.Del(ctx, ks...).Err()
					}
				} else {
					delErr = cli.Del(ctx, ks...).Err()
				}

				mutex.Lock()
				deletedKeys = append(deletedKeys, ks...)
				mutex.Unlock()

				return delErr
			}

			scanSize := int64(Preferences().GetScanSize())
			iter := cli.Scan(ctx, 0, key, scanSize).Iterator()
			resultKeys := make([]string, 0, 100)
			for iter.Next(ctx) {
				resultKeys = append(resultKeys, iter.Val())
				if len(resultKeys) >= 20 {
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
			if err = client.Unlink(ctx, key).Err(); err != nil {
				if err = client.Del(ctx, key).Err(); err != nil {
					resp.Msg = err.Error()
					return
				}
			}
		} else {
			if err = client.Del(ctx, key).Err(); err != nil {
				resp.Msg = err.Error()
				return
			}
		}
		deletedKeys = append(deletedKeys, key)
	}

	resp.Success = true
	resp.Data = map[string]any{
		"deleted":     deletedKeys,
		"deleteCount": len(deletedKeys),
	}
	return
}

// FlushDB flush database
func (b *browserService) FlushDB(connName string, db int, async bool) (resp types.JSResp) {
	item, err := b.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	flush := func(ctx context.Context, cli redis.UniversalClient, async bool) error {
		_, e := cli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Select(ctx, db)
			if async {
				pipe.FlushDBAsync(ctx)
			} else {
				pipe.FlushDB(ctx)
			}
			return nil
		})
		return e
	}

	client, ctx := item.client, item.ctx
	if cluster, ok := client.(*redis.ClusterClient); ok {
		// cluster mode
		err = cluster.ForEachMaster(ctx, func(ctx context.Context, cli *redis.Client) error {
			return flush(ctx, cli, async)
		})
		// try sync mode if error cause
		if err != nil && async {
			err = cluster.ForEachMaster(ctx, func(ctx context.Context, cli *redis.Client) error {
				return flush(ctx, cli, false)
			})
		}
	} else {
		if err = flush(ctx, client, async); err != nil && async {
			// try sync mode if error cause
			err = flush(ctx, client, false)
		}
	}

	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Success = true
	return
}

// RenameKey rename key
func (b *browserService) RenameKey(connName string, db int, key, newKey string) (resp types.JSResp) {
	item, err := b.getRedisClient(connName, db)
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
func (b *browserService) GetCmdHistory(pageNo, pageSize int) (resp types.JSResp) {
	resp.Success = true
	if pageSize <= 0 || pageNo <= 0 {
		// return all history
		resp.Data = map[string]any{
			"list":     b.cmdHistory,
			"pageNo":   1,
			"pageSize": -1,
		}
	} else {
		total := len(b.cmdHistory)
		startIndex := total / pageSize * (pageNo - 1)
		endIndex := min(startIndex+pageSize, total)
		resp.Data = map[string]any{
			"list":     b.cmdHistory[startIndex:endIndex],
			"pageNo":   pageNo,
			"pageSize": pageSize,
		}
	}
	return
}

// CleanCmdHistory clean redis command history
func (b *browserService) CleanCmdHistory() (resp types.JSResp) {
	b.cmdHistory = []cmdHistoryItem{}
	resp.Success = true
	return
}

// GetSlowLogs get slow log list
func (b *browserService) GetSlowLogs(connName string, db int, num int64) (resp types.JSResp) {
	item, err := b.getRedisClient(connName, db)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	client, ctx := item.client, item.ctx
	var logs []redis.SlowLog
	if cluster, ok := client.(*redis.ClusterClient); ok {
		// cluster mode
		var mu sync.Mutex
		err = cluster.ForEachShard(ctx, func(ctx context.Context, cli *redis.Client) error {
			if subLogs, _ := client.SlowLogGet(ctx, num).Result(); len(subLogs) > 0 {
				mu.Lock()
				logs = append(logs, subLogs...)
				mu.Unlock()
			}
			return nil
		})
	} else {
		logs, err = client.SlowLogGet(ctx, num).Result()
	}
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	sort.Slice(logs, func(i, j int) bool {
		return logs[i].Time.UnixMilli() > logs[j].Time.UnixMilli()
	})
	if len(logs) > int(num) {
		logs = logs[:num]
	}

	list := sliceutil.Map(logs, func(i int) slowLogItem {
		var name string
		var e error
		if name, e = url.QueryUnescape(logs[i].ClientName); e != nil {
			name = logs[i].ClientName
		}
		return slowLogItem{
			Timestamp: logs[i].Time.UnixMilli(),
			Client:    name,
			Addr:      logs[i].ClientAddr,
			Cmd:       sliceutil.JoinString(logs[i].Args, " "),
			Cost:      logs[i].Duration.Milliseconds(),
		}
	})

	resp.Success = true
	resp.Data = map[string]any{
		"list": list,
	}
	return
}
