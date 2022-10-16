package fuse

import "github.com/phuhao00/network"

type Principal interface {
	Resolve(*network.Packet)
}
