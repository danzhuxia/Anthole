package client

import (
	"Anthole/pkg/common"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var mutex sync.Mutex

type TcpWorker struct {
	serverConn net.Conn
	verify     bool
}

func (t *TcpWorker) Run(config *common.AntHoleConfig) error {

	for {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port))
		if err != nil {
			logrus.Printf("与服务端连接失败，10秒后重试,错误信息： %v", err)

			time.Sleep(10 * time.Second)
			continue
		}
		t.serverConn = conn
		// tcp连接成功后发送权限包 验证权限

		t.SendDataToServer(common.Transmission{
			RequestId:  common.VerifyKey,
			ServiceId:  0,
			Data:       []byte(config.Common.Token),
			DataLength: uint16(len([]byte(config.Common.Token))),
		})

		t.serverDataHandle(config, conn)
		logrus.Info("与服务端连接断开，10秒后重试")
		time.Sleep(10 * time.Second)
	}
}

func (t *TcpWorker) SendDataToServer(data common.Transmission) {
	logrus.Debug("发送数据给Server")
	mutex.Lock()
	defer mutex.Unlock()
	t.serverConn.Write(data.ConvertToByte())
}

func (t *TcpWorker) serverDataHandle(config *common.AntHoleConfig, conn net.Conn) {
	for {
		//读取 conn 中的数据
		buf, err := common.GetDataFromConn(conn, common.TransmissionPackageLength)
		if err != nil {
			logrus.Info("与服务端连接断开")
			conn.Close()
			break
		}

		pkg := common.TodoTransmission(buf)
		logrus.Debugf("读取到服务端数据包，Requestid %d ， ServiceID: %d , DataLength:%d ", pkg.RequestId, pkg.ServiceId, pkg.DataLength)

		if pkg.RequestId == common.CloseSocket {
			// 处理系统消息
			data := string(pkg.GetData())
			logrus.Debug("关闭本地socket", data)

			info := strings.Split(data, ",")

			poolKey := fmt.Sprintf("%s_%s", info[1], info[0])
			pool_conn, ok := antholeClientConnPool.Load(poolKey)

			if ok {
				(pool_conn.(net.Conn)).Close()
				antholeClientConnPool.Delete(poolKey)
			}

			continue
		}

		if pkg.RequestId == common.OpenSocket {
			// 处理系统消息
			data := string(pkg.GetData())

			info := strings.Split(data, ",") // requestid , serviceid

			requestId, _ := strconv.ParseUint(info[0], 10, 64)
			serviceId, _ := strconv.ParseUint(info[1], 10, 16)
			poolKey := fmt.Sprintf("%s_%s", info[1], info[0])
			_, ok := antholeClientConnPool.Load(poolKey)

			// 已有就不处理了
			if ok {
				continue
			} else {
				createLocalSocket(config, requestId, uint16(serviceId), poolKey)
			}
			continue
		}
		logrus.Debug("接收到Server数据")
		//log.Println(string(pkg.Data))

		poolKey := fmt.Sprintf("%d_%d", pkg.ServiceId, pkg.RequestId)

		pool_conn, ok := antholeClientConnPool.Load(poolKey)
		if !ok {
			// 没有本地连接的时候先创建连接
			logrus.Debug("创建位置2")
			pool_conn = createLocalSocket(config, pkg.RequestId, pkg.ServiceId, poolKey)
		}
		_, err = (pool_conn.(net.Conn)).Write(pkg.GetData())
		if err != nil {
			logrus.Debug("写入本地socket错误 %v", err)
			(pool_conn.(net.Conn)).Close()
			antholeClientConnPool.Delete(poolKey)
		}
	}
}
