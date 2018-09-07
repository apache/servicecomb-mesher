package grpc

import (
	"fmt"
	common2 "github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-mesh/mesher/common"
	"github.com/go-mesh/mesher/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	authority    = ":authority"
	MetadataPort = "forwarded-port"
)

var dr = resolver.GetDestinationResolver(protocolName)

//LocalRequestHandler handle local request
func LocalRequestHandler(srv interface{}, serverStream grpc.ServerStream) error {
	// what we can get from grpc context
	// context.Background.WithCancel.
	// WithDeadline(2018-08-10 14:57:40.520799844 +0800 CST m=+1372.256775226 [998.578561ms]).
	// WithValue(peer.peerKey{}, &peer.Peer{Addr:(*net.TCPAddr)(0xc42026ede0), AuthInfo:credentials.AuthInfo(nil)}).
	// WithValue(metadata.mdIncomingKey{}, metadata.MD{":authority":[]string{"localhost:40101"}, "content-type":[]string{"application/grpc"}, "user-agent":[]string{"grpc-go/1.14.0-dev"}}).
	// WithValue(grpc.streamKey{}, <stream: 0xc42026a100, /helloworld.Greeter/SayHello>)

	inv, err := transfer2Invocation(serverStream, true)
	if err != nil {
		return err
	}
	c, err := handler.GetChain(common2.Consumer, common.ChainConsumerOutgoing)
	if err != nil {
		lager.Logger.Error("Get chain failed: " + err.Error())
		return err
	}
	var invRsp *invocation.Response
	c.Next(inv, func(ir *invocation.Response) error {
		//Send the request to the destination
		invRsp = ir
		if invRsp != nil {
			return invRsp.Err
		}
		return nil
	})
	if invRsp.Err != nil {
		return invRsp.Err
	}
	err = serverStream.SendMsg(invRsp.Result)
	if err != nil {
		return err
	}
	return nil
}

//RemoteRequestHandler handle remote request
func RemoteRequestHandler(srv interface{}, serverStream grpc.ServerStream) error {
	inv, err := transfer2Invocation(serverStream, false)
	if err != nil {
		return err
	}
	c, err := handler.GetChain(common2.Provider, common.ChainProviderIncoming)
	if err != nil {
		lager.Logger.Error("Get chain failed: " + err.Error())
		return err
	}
	var invRsp *invocation.Response
	c.Next(inv, func(ir *invocation.Response) error {
		//Send the request to the destination
		invRsp = ir
		if invRsp != nil {
			return invRsp.Err
		}
		return nil
	})
	if invRsp.Err != nil {
		return invRsp.Err
	}
	err = serverStream.SendMsg(invRsp.Result)
	if err != nil {
		return err
	}
	return nil
}

//prepare headers and service name
// it assigns stream ctx to invocation ctx
func transfer2Invocation(stream grpc.ServerStream, FromLocal bool) (*invocation.Invocation, error) {
	var err error
	inv := invocation.New(stream.Context())
	h := inv.Headers()
	//we can get service name and header is here
	md, _ := metadata.FromIncomingContext(stream.Context())
	for k := range md {
		h[k] = md.Get(k)[0]
	}
	// schema and operation is in stream context
	schemaAndOperation, ok := grpc.MethodFromServerStream(stream)
	if !ok {
		return nil, fmt.Errorf("can not full method name")
	}
	p, ok := peer.FromContext(stream.Context())
	if !ok {
		return nil, fmt.Errorf("can not get peer")
	}
	inv.Protocol = "grpc"
	inv.OperationID = schemaAndOperation
	inv.Args = stream
	inv.Reply = &frame{}
	if FromLocal {
		h[MetadataPort], err = dr.Resolve(p.Addr.String(), h, h[authority], &inv.MicroServiceName)
		if err != nil {
			return nil, err
		}
	} else {
		inv.Endpoint = "127.0.0.1:" + h[MetadataPort]
	}
	return inv, nil
}
