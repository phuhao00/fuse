package fuse

import (
	"errors"
	"fmt"
	"github.com/phuhao00/network"
)

var (
	// ErrNotFound is returned when no route match is found.
	ErrNotFound = errors.New("no matching route was found")
)

func NewRouter() *Router {
	return &Router{namedRoutes: make(map[string]*Route)}
}

type Router struct {
	NotFoundHandler Handler

	MethodNotAllowedHandler Handler

	routes []*Route

	namedRoutes map[string]*Route

	KeepContext bool

	middlewares []middleware

	routeConf
}

type routeConf struct {
	matchers []matcher
}

// RouteMatch stores information about a matched route.
type RouteMatch struct {
	Handler Handler
	Vars    map[string]string

	MatchErr error
}

func (r *Router) NewRoute() *Route {
	route := &Route{namedRoutes: r.namedRoutes}
	r.routes = append(r.routes, route)
	return route
}

func (r *Router) Match(req *network.Packet, match *RouteMatch) bool {
	for _, route := range r.routes {
		if route.Match(req, match) {
			for i := len(r.middlewares) - 1; i >= 0; i-- {
				match.Handler = r.middlewares[i].Middleware(match.Handler)
			}
			return true
		}
	}
	return false
}

func (r *Router) GetHandler(req *network.Packet) Handler {

	var match RouteMatch

	if !r.Match(req, &match) {
		return nil
	}

	return match.Handler

}

func (r *Router) AddRoute(Id uint64, f Handler) *Route {
	return r.NewRoute().SetId(Id).SetHandler(f)
}

func (r *Route) SetHandler(handler Handler) *Route {
	if r.err == nil {
		r.handler = handler
	}
	return r
}

func (r *Route) GetHandler() Handler {
	return r.handler
}

func (r *Route) SetDescription(desc string) *Route {
	if r.Description != "" {
		r.err = fmt.Errorf("route already has Description %q, can't set %q",
			r.Description, desc)
	}
	if r.err == nil {
		r.Description = desc
	}
	return r
}

func (r *Route) GetName() string {
	return r.Description
}

func (r *Route) GetError() error {
	return r.err
}

func (r *Route) Match(packet *network.Packet, match *RouteMatch) bool {
	if r.err != nil {
		return false
	}

	var matchErr error

	for _, m := range r.matchers {
		if matched := m.Match(packet, match); !matched {
			matchErr = ErrNotFound
			continue
		} else {
			matchErr = nil
			break
		}
	}

	if matchErr != nil {
		match.MatchErr = matchErr
		return false
	}

	match.Handler = r.handler

	return true
}
