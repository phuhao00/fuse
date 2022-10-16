package fuse

import (
	"fmt"
	"github.com/phuhao00/network"
	"testing"
)

type PrincipalDemo struct {
}

func (p *PrincipalDemo) Resolve(packet *network.Packet) {
	fmt.Println("impl")
}

type Manager struct {
	Pr map[uint64]*PrincipalDemo
	Ch chan *network.Packet
}

var (
	testRouter *Router
)

func (m *Manager) Run() {
	for {
		select {
		case req := <-m.Ch:
			handler := testRouter.GetHandler(req)
			pr := m.GetPrincipal()
			handler(req, pr)
		}
	}
}

func (m *Manager) GetPrincipal() Principal {
	return m.Pr[22]
}

func TestFuse(t *testing.T) {
	var (
		m *Manager
	)
	testRouter = NewRouter()
	m = &Manager{}
	p := &PrincipalDemo{}
	m.Pr[22] = p
	testRouter.AddRoute(11, func(packet *network.Packet, p Principal) {
		p.Resolve(packet)
	})
	m.Run()
}
