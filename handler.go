package fuse

import (
	"github.com/phuhao00/network"
	_ "github.com/phuhao00/network"
)

type Handler func(packet *network.Packet, p Principal)
