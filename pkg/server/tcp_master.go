package server

import (
	"Anthole/pkg/common"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var mutex sync.Mutex

type TcpMaster struct {
	antholeClientConn net.Conn
	verify            bool
}

func (t *TcpMaster) Run(port int) error {
	listen, err := net.Listen("tcp", strconv.Itoa(port))
	if err != nil {
		logrus.Errorf("error listen:%v", err)
		return err
	}
	logrus.Info("Anthole 启动成功！")
	defer listen.Close()

	for {
		if conn, err := listen.Accept(); err != nil {
			logrus.Errorf("accept error:%v", err)

		} else {
			logrus.Info("客户端接入")
			// 设置一个3秒的计时器，3秒后没通过权限验证则关闭连接
			go func() {
				time.Sleep(3 * time.Second)
				if !t.verify {
					logrus.Info("权限验证超时，关闭连接")
					conn.Close()
				}
			}()
			t.antholeClientConn = conn
			t.ClientDataHandle(conn)
		}
	}
}

func (t *TcpMaster) ClientDataHandle(conn net.Conn) {
	//循环读取RequestClient数据流
	for {
		//网络数据流读入 buffer
		buf, err := common.GetDataFromConn(conn, common.TransmissionPackageLength)
		//数据读尽、读取错误 关闭 socket 连接
		if err != nil {
			logrus.Infof("客户端连接关闭,%v", err)
			t.verify = false
			conn.Close()
			break
		}

		dataPackage := common.TodoTransmission(buf)

		if !t.verify && dataPackage.RequestId != common.VerifyKey {
			// 验证未通过切不是权限验证的数据包
			conn.Close()
			t.verify = false
			logrus.Info("权限校验失败！关闭连接")
			break
		}

		if dataPackage.RequestId == common.VerifyKey {
			config, _ := common.GetConfig("")
			if string(dataPackage.GetData()) == config.Common.Token {
				logrus.Info("权限校验通过")
				t.verify = true
			} else {
				conn.Close()
				t.verify = false
				logrus.Info("权限校验失败！关闭连接")
				break
			}
		}
		if dataPackage.RequestId == common.CloseSocket {
			// 处理系统消息
			data := string(dataPackage.GetData())
			logrus.Info("关闭Request Socket", data)
			info := strings.Split(data, ",") // request_id  serviceid
			poolKey := fmt.Sprintf("%v_%v", info[1], info[0])
			CloseRequestSocket(poolKey)
			continue
		}
		// 将数据包发送给request
		sendDataToRequest(dataPackage)
	}
}

func (t *TcpMaster) SendDataToClient(data common.Transmission) (int, error) {
	mutex.Lock()
	defer mutex.Unlock()
	if t.antholeClientConn != nil && t.verify {
		encrData := data.ConvertToByte()
		logrus.Debugf("发送数据到Client,加密包长度%d ,RequestID: %d, ServiceID: %d,DataLength: %d  ",
			len(encrData), data.RequestId, data.ServiceId, data.DataLength)
		return t.antholeClientConn.Write(encrData)
	}
	return 0, errors.New("客户端未连接")
}
