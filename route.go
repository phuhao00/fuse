package fuse

import (
	"github.com/phuhao00/network"
)

type Route struct {
	Id uint64
	// Request handler for the route.
	handler Handler
	// If true, this route never matches: it is only used to build URLs.
	buildOnly bool
	// The Description used to build URLs.
	Description string
	// Error resulted from building a route.
	err error

	// "global" reference to all named routes
	namedRoutes map[string]*Route

	// config possibly passed in from `Router`
	routeConf
}

type matcher interface {
	Match(*network.Packet, *RouteMatch) bool
}

func (r *Route) SetId(Id uint64) *Route {
	r.Id = Id
	return r
}
