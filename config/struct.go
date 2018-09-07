package config

//MesherConfig has all mesher config
type MesherConfig struct {
	PProf       *PProf         `yaml:"pprof"`
	Plugin      *Plugin        `yaml:"plugin"`
	Admin       Admin          `yaml:"admin"`
	HealthCheck []*HealthCheck `yaml:"localHealthCheck"`
	ProxyedPro  string         `yaml:"proxyedProtocol"`
}

//HealthCheck define how to check local ports
type HealthCheck struct {
	Port     string `yaml:"port"`
	Protocol string `yaml:"protocol"`
	URI      string `yaml:"uri"`
	Interval string `yaml:"interval"`
	Match    *Match `yaml:"match"`
}

//Match define health check result success criteria
type Match struct {
	Status string `yaml:"status"`
	Body   string `yaml:"body"`
}

//PProf has enable and listen attribute for pprof
type PProf struct {
	Enable bool   `yaml:"enable"`
	Listen string `yaml:"listen"`
}

//Policy has attributes for destination, tags and loadbalance
type Policy struct {
	Destination   string            `yaml:"destination"`
	Tags          map[string]string `yaml:"tags"`
	LoadBalancing map[string]string `yaml:"loadBalancing"`
}

//Plugin has attributes for destination and source resolver
type Plugin struct {
	DestinationResolver map[string]string `yaml:"destinationResolver"`
	SourceResolver      string            `yaml:"sourceResolver"`
}

//Admin has attributes for enabling, serverURI and metrics for admin data
type Admin struct {
	Enable           *bool  `yaml:"enable"`
	ServerURI        string `yaml:"serverUri"`
	GoRuntimeMetrics bool   `yaml:"goRuntimeMetrics"`
}
