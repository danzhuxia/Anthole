package client

import (
	"Anthole/pkg/common"
	"fmt"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
)

//本地Anthole连接池
var antholeClientConnPool sync.Map

type ClientWorker interface {
	Run(config *common.AntHoleConfig) error
	SendDataToServer(data common.Transmission)
	serverDataHandle(config *common.AntHoleConfig, conn net.Conn)
}

var mu sync.Mutex
var instance ClientWorker

func StartClientWorkerInstance() ClientWorker {
	mu.Lock()
	defer mu.Unlock()
	if instance == nil {
		instance = &TcpWorker{
			verify: true,
		}
	}
	return instance
}

// 读取本地socket数据
func connHandleLocal(requestId uint64, serviceId uint16, conn net.Conn) {

	//缓存 conn 中的数据
	buf := make([]byte, common.TransmissionDataLength)

	for {

		cnt, err := conn.Read(buf)

		if cnt == 0 || err != nil {
			conn.Close()
			poolKey := fmt.Sprintf("%d_%d", serviceId, requestId)
			antholeClientConnPool.Delete(poolKey)

			data := []byte(fmt.Sprintf("%d,%d", requestId, serviceId))

			logrus.Debug("通知服务端断开socket %d %d", requestId, serviceId)
			// 通知服务端断开socket
			StartClientWorkerInstance().SendDataToServer(common.Transmission{
				RequestId:  common.CloseSocket,
				ServiceId:  serviceId,
				Data:       data,
				DataLength: uint16(len(data)),
			})

			break
		}

		pkg := common.Transmission{
			RequestId:  requestId,
			ServiceId:  serviceId,
			DataLength: uint16(cnt),
			Data:       buf,
		}

		logrus.Debug("收到本地Socket数据,长度%d", cnt)

		// 发送给server
		StartClientWorkerInstance().SendDataToServer(pkg)

	}
}

func createLocalSocket(config *common.AntHoleConfig, requestId uint64, serverId uint16, poolKey string) net.Conn {
	// 没有连接的时候先创建连接
	local_conn, err := net.Dial("tcp",
		fmt.Sprintf("%s:%d", config.Client.Services[serverId].LocalHost, config.Client.Services[serverId].LocalPort))
	if err == nil {
		logrus.Debug("创建本地连接成功")
		antholeClientConnPool.Store(poolKey, local_conn)
		go connHandleLocal(requestId, serverId, local_conn)
	} else {
		logrus.Fatalf("本地连接错误:%v", err)
	}

	return local_conn
}
