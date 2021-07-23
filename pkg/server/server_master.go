package server

import (
	"Anthole/pkg/common"
	"net"
	"sync"
)

type Master interface {
	Run(int) error
	ClientDataHandle(net.Conn)
	SendDataToClient(common.Transmission) (int, error)
}

var instance Master
var mu sync.Mutex

func StartMasterInstance() Master {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		instance = &TcpMaster{
			verify: false,
		}
	}
	return instance
}
