package services

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/klauspost/compress/zip"
	"github.com/redis/go-redis/v9"
	"github.com/vrischmann/userdir"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"path"
	"strings"
	"sync"
	"time"
	. "tinyrdm/backend/storage"
	"tinyrdm/backend/types"
)

type cmdHistoryItem struct {
	Timestamp int64  `json:"timestamp"`
	Server    string `json:"server"`
	Cmd       string `json:"cmd"`
	Cost      int64  `json:"cost"`
}

type connectionService struct {
	ctx   context.Context
	conns *ConnectionsStorage
}

var connection *connectionService
var onceConnection sync.Once

func Connection() *connectionService {
	if connection == nil {
		onceConnection.Do(func() {
			connection = &connectionService{
				conns: NewConnections(),
			}
		})
	}
	return connection
}

func (c *connectionService) Start(ctx context.Context) {
	c.ctx = ctx
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

		tlsConfig = &tls.Config{
			RootCAs:            caCertPool,
			InsecureSkipVerify: config.SSL.AllowInsecure,
			Certificates:       certs,
			ServerName:         strings.TrimSpace(config.SSL.SNI),
		}
	}

	option := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Addr, config.Port),
		Username:     config.Username,
		Password:     config.Password,
		DialTimeout:  time.Duration(config.ConnTimeout) * time.Second,
		ReadTimeout:  time.Duration(config.ExecTimeout) * time.Second,
		WriteTimeout: time.Duration(config.ExecTimeout) * time.Second,
		TLSConfig:    tlsConfig,
	}
	if config.LastDB > 0 {
		option.DB = config.LastDB
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

// SaveLastDB save last selected database index
func (c *connectionService) SaveLastDB(name string, db int) (resp types.JSResp) {
	param := c.conns.GetConnection(name)
	if param == nil {
		resp.Msg = "no connection named \"" + name + "\""
		return
	}

	if param.LastDB != db {
		param.LastDB = db
		if err := c.conns.UpdateConnection(name, param.ConnectionConfig); err != nil {
			resp.Msg = "save connection fail:" + err.Error()
			return
		}
	}
	resp.Success = true
	return
}

// SaveRefreshInterval save auto refresh interval
func (c *connectionService) SaveRefreshInterval(name string, interval int) (resp types.JSResp) {
	param := c.conns.GetConnection(name)
	if param == nil {
		resp.Msg = "no connection named \"" + name + "\""
		return
	}
	if param.RefreshInterval != interval {
		param.RefreshInterval = interval
		if err := c.conns.UpdateConnection(name, param.ConnectionConfig); err != nil {
			resp.Msg = "save connection fail:" + err.Error()
			return
		}
	}
	resp.Success = true
	return
}

// ExportConnections export connections to zip file
func (c *connectionService) ExportConnections() (resp types.JSResp) {
	defaultFileName := "connections_" + time.Now().Format("20060102150405") + ".zip"
	filepath, err := runtime.SaveFileDialog(c.ctx, runtime.SaveDialogOptions{
		ShowHiddenFiles: true,
		DefaultFilename: defaultFileName,
		Filters: []runtime.FileFilter{
			{
				Pattern: "*.zip",
			},
		},
	})
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	// compress the connections profile with zip
	const connectionFilename = "connections.yaml"
	inputFile, err := os.Open(path.Join(userdir.GetConfigHome(), "TinyRDM", connectionFilename))
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	defer inputFile.Close()

	outputFile, err := os.Create(filepath)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	defer outputFile.Close()

	zipWriter := zip.NewWriter(outputFile)
	defer zipWriter.Close()

	headerWriter, err := zipWriter.CreateHeader(&zip.FileHeader{
		Name:   connectionFilename,
		Method: zip.Deflate,
	})
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	if _, err = io.Copy(headerWriter, inputFile); err != nil {
		resp.Msg = err.Error()
		return
	}

	resp.Success = true
	resp.Data = struct {
		Path string `json:"path"`
	}{
		Path: filepath,
	}
	return
}

// ImportConnections import connections from local zip file
func (c *connectionService) ImportConnections() (resp types.JSResp) {
	filepath, err := runtime.OpenFileDialog(c.ctx, runtime.OpenDialogOptions{
		ShowHiddenFiles: true,
		Filters: []runtime.FileFilter{
			{
				Pattern: "*.zip",
			},
		},
	})
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	const connectionFilename = "connections.yaml"
	zipFile, err := zip.OpenReader(filepath)
	if err != nil {
		resp.Msg = err.Error()
		return
	}

	var file *zip.File
	for _, file = range zipFile.File {
		if file.Name == connectionFilename {
			break
		}
	}
	if file != nil {
		zippedFile, err := file.Open()
		if err != nil {
			resp.Msg = err.Error()
			return
		}
		defer zippedFile.Close()

		outputFile, err := os.Create(path.Join(userdir.GetConfigHome(), "TinyRDM", connectionFilename))
		if err != nil {
			resp.Msg = err.Error()
			return
		}
		defer outputFile.Close()

		if _, err = io.Copy(outputFile, zippedFile); err != nil {
			resp.Msg = err.Error()
			return
		}
	}

	resp.Success = true
	return
}
