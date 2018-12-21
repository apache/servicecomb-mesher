package grpc

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/go-chassis/go-chassis/core/client"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
	"golang.org/x/net/http2"
	"net"
	"net/http"
)

const (
	//SchemaHTTP represents the http schema
	SchemaHTTP = "http"
	//SchemaHTTPS represents the https schema
	SchemaHTTPS = "https"
)

//ErrInvalidResp invalid input
var (
	ErrInvalidResp = errors.New("rest consumer response arg is not *http.Response type")
	//ErrCanceled means Request is canceled by context management
	ErrCanceled = errors.New("request cancelled")
)

//Client is a grpc client
type Client struct {
	c    *http.Client
	opts client.Options
}

var protocolName = "grpc"

func init() {
	client.InstallPlugin(protocolName, NewClient)
}

//NewClient return a new client of grpc
func NewClient(opts client.Options) (client.ProtocolClient, error) {
	client := &http.Client{}
	if opts.TLSConfig != nil {
		client.Transport = &http2.Transport{
			TLSClientConfig: opts.TLSConfig,
		}
	} else {
		client.Transport = &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			}}
	}
	return &Client{
		c:    client,
		opts: opts,
	}, nil
}
func (c *Client) contextToHeader(ctx context.Context, req *http.Request) {
	for k, v := range common.FromContext(ctx) {
		req.Header.Set(k, v)
	}
}

//Call is a method which uses grpc protocol to transfer invocation
func (c *Client) Call(ctx context.Context, addr string, inv *invocation.Invocation, rsp interface{}) error {
	var err error
	reqSend, err := httputil.HTTPRequest(inv)
	if err != nil {
		return err
	}
	resp, ok := rsp.(*http.Response)
	if !ok {
		return ErrInvalidResp
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if c.opts.TLSConfig != nil {
		reqSend.URL.Scheme = SchemaHTTPS
	} else {
		reqSend.URL.Scheme = SchemaHTTP
	}
	if addr != "" {
		reqSend.URL.Host = addr
	}

	var temp *http.Response
	errChan := make(chan error, 1)
	go func() {
		temp, err = c.c.Do(reqSend)
		errChan <- err
	}()

	select {
	case <-ctx.Done():
		err = ErrCanceled
	case err = <-errChan:
		if err == nil {
			*resp = *temp
		}
	}
	return err
}

//String return name
func (c *Client) String() string {
	return protocolName
}

//Close close the conn
func (c *Client) Close() error {
	return nil
}
