package ports

import "github.com/go-chassis/go-chassis/core/common"

var defaultProtocolPort = map[string]string{
	common.ProtocolRest: "30101",
	"grpc":              "40101",
}

//SetFixedPort allows developer set a fixed port for for you protocol
func SetFixedPort(protocol, port string) {
	defaultProtocolPort[protocol] = port
}

//GetFixedPort return port pf a protocol
func GetFixedPort(protocol string) string {
	return defaultProtocolPort[protocol]

}
