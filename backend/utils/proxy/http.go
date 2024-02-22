package proxy

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

type HttpProxy struct {
	scheme  string       // HTTP Proxy scheme
	host    string       // HTTP Proxy host or host:port
	auth    *proxy.Auth  // authentication
	forward proxy.Dialer // forwarding Dialer
}

func (p *HttpProxy) Dial(network, addr string) (net.Conn, error) {
	c, err := p.forward.Dial(network, p.host)
	if err != nil {
		return nil, err
	}

	err = c.SetDeadline(time.Now().Add(15 * time.Second))
	if err != nil {
		return nil, err
	}

	reqUrl := &url.URL{
		Scheme: "",
		Host:   addr,
	}

	// create with CONNECT method
	req, err := http.NewRequest("CONNECT", reqUrl.String(), nil)
	if err != nil {
		c.Close()
		return nil, err
	}
	req.Close = false

	// authentication
	if p.auth != nil {
		req.SetBasicAuth(p.auth.User, p.auth.Password)
		req.Header.Add("Proxy-Authorization", req.Header.Get("Authorization"))
	}

	// send request
	err = req.Write(c)
	if err != nil {
		c.Close()
		return nil, err
	}

	res, err := http.ReadResponse(bufio.NewReader(c), req)
	if err != nil {
		res.Body.Close()
		c.Close()
		return nil, err
	}
	res.Body.Close()

	if res.StatusCode != http.StatusOK {
		c.Close()
		return nil, fmt.Errorf("proxy connection error: StatusCode[%d]", res.StatusCode)
	}

	return c, nil
}

func NewHttpProxyDialer(u *url.URL, forward proxy.Dialer) (proxy.Dialer, error) {
	var auth *proxy.Auth
	if u.User != nil {
		pwd, _ := u.User.Password()
		auth = &proxy.Auth{
			User:     u.User.Username(),
			Password: pwd,
		}
	}

	hp := &HttpProxy{
		scheme:  u.Scheme,
		host:    u.Host,
		auth:    auth,
		forward: forward,
	}
	return hp, nil
}

func init() {
	proxy.RegisterDialerType("http", NewHttpProxyDialer)
	proxy.RegisterDialerType("https", NewHttpProxyDialer)
}
