package pilot

import (
	"fmt"
	"github.com/go-mesh/mesher/proxy/config"
	"github.com/go-mesh/mesher/proxy/pkg/egress"
	"sync"
)

func init() { egress.InstallEgressService("pilot", newPilotEgress) }

func newPilotEgress() (egress.Egress, error) { return &PilotEgress{}, nil }

//PilotEgress is pilot egress service
type PilotEgress struct{}

//FetchEgressRule return all rules
func (r *PilotEgress) FetchEgressRule() map[string][]*config.EgressRule {
	return GetEgressRule()
}

//SetEgressRule set rules
func (r *PilotEgress) SetEgressRule(rr map[string][]*config.EgressRule) {
	SetEgressRule(rr)
}

//Init init egress config
func (r *PilotEgress) Init(o egress.Options) error {
	// the manager use dests to init, so must init after dests
	if err := InitPilotFetcher(o); err != nil {
		return err
	}
	return refresh()
}

// refresh all the egress config
func refresh() error {
	configs := pilotfetcher.GetConfigurations()

	d := make(map[string][]*config.EgressRule)
	for k, v := range configs {
		rules, ok := v.([]*config.EgressRule)
		if !ok {
			err := fmt.Errorf("Egress rule type assertion fail, key: %s", k)
			return err
		}
		d[k] = rules
	}
	ok, _ := egress.ValidateEgressRule(d)
	if ok {
		dests = d
	}
	return nil
}

var dests = make(map[string][]*config.EgressRule)
var lock sync.RWMutex

// GetEgressRule get egress rule
func GetEgressRule() map[string][]*config.EgressRule {
	lock.RLock()
	defer lock.RUnlock()
	return dests
}

// SetEgressRule set egress rule
func SetEgressRule(rule map[string][]*config.EgressRule) {
	lock.RLock()
	defer lock.RUnlock()
	dests = rule
}
