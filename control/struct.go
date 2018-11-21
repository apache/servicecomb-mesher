package control

//EgressConfig is a standardized model
type EgressConfig struct {
	Hosts []string
	Ports []*EgressPort
}

//EgressPort protocol and the corresponding port
type EgressPort struct {
	Port     int32
	Protocol string
}
