package grpc

import (
	"context"
	"github.com/go-chassis/go-chassis/core/client"
	"github.com/go-chassis/go-chassis/core/invocation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

//Client is a grpc client
type Client struct {
	conn *grpc.ClientConn
}

func init() {
	client.InstallPlugin(protocolName, NewClient)
}

var sd = &grpc.StreamDesc{
	ClientStreams: true,
	ServerStreams: true,
}

//NewClient return a new client of grpc
func NewClient(opts client.Options) (client.ProtocolClient, error) {
	var cc *grpc.ClientConn
	var err error
	if opts.TLSConfig == nil {
		cc, err = grpc.Dial(opts.Endpoint, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
	} else {
		cred := credentials.NewTLS(opts.TLSConfig)
		cc, err = grpc.Dial(opts.Endpoint, grpc.WithTransportCredentials(cred))
	}

	return &Client{
		conn: cc,
	}, nil
}

//Call is a method which uses grpc protocol to transfer invocation
func (c *Client) Call(ctx context.Context, addr string, inv *invocation.Invocation, rsp interface{}) error {
	serverStream, _ := inv.Args.(grpc.ServerStream)
	//send headers to provider, including incoming md
	md := metadata.MD{}
	for k, v := range inv.Headers() {
		md.Set(k, v)
	}
	ctx = metadata.NewOutgoingContext(serverStream.Context(), md)

	clientStream, err := grpc.NewClientStream(ctx, sd, c.conn, inv.OperationID, grpc.CallCustomCodec(&codec{}))
	if err != nil {
		return err
	}
	f := rsp.(*frame)
	//get frame
	err = serverStream.RecvMsg(f)
	if err != nil {
		return err
	}
	//set frame to provider
	err = clientStream.SendMsg(f)
	if err != nil {
		return err
	}
	//receive frame from provider
	err = clientStream.RecvMsg(f)
	if err != nil {
		return err
	}

	rsp = f

	return nil
}

//String return name
func (c *Client) String() string {
	return protocolName
}

//Close close the conn
func (c *Client) Close() error {
	return c.conn.Close()
}
